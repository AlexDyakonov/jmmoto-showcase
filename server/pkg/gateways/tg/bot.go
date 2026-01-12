package tg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shampsdev/go-telegram-template/pkg/config"
	"github.com/shampsdev/go-telegram-template/pkg/domain"
	"github.com/shampsdev/go-telegram-template/pkg/repo"
	"github.com/shampsdev/go-telegram-template/pkg/usecase"
	"github.com/shampsdev/go-telegram-template/pkg/utils/slogx"
)

type Bot struct {
	*bot.Bot
	cases usecase.Cases
	log   *slog.Logger

	botUrl    string
	webAppUrl string

	// Состояние ожидания ввода цены: userID -> motorcycleID
	waitingPrice sync.Map
	// Состояние ожидания ввода даты прибытия: userID -> motorcycleID
	waitingArrivalDate sync.Map
}

func NewBot(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) (*Bot, error) {
	opts := []bot.Option{}

	if cfg.Debug {
		opts = append(opts, bot.WithDebug())
	}
	tgb, err := bot.New(cfg.TG.BotToken, opts...)
	if err != nil {
		return nil, fmt.Errorf("error creating bot: %w", err)
	}
	cases := usecase.Setup(ctx, cfg, pool)

	b := &Bot{
		Bot:   tgb,
		cases: cases,
		log:   slogx.FromCtx(ctx),
	}

	me, err := b.GetMe(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting bot info: %w", err)
	}
	b.webAppUrl = fmt.Sprintf("https://t.me/%s/%s", me.Username, cfg.TG.WebAppName)
	b.botUrl = fmt.Sprintf("https://t.me/%s", me.Username)

	return b, nil
}

func (b *Bot) Run(ctx context.Context) {
	_, err := b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{},
	})
	if err != nil {
		panic(fmt.Errorf("error setting bot commands: %w", err))
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, b.handleCommandStart)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypePrefix, b.handleMessage)

	b.Start(ctx)
}

func (b *Bot) handleCommandStart(ctx context.Context, _ *bot.Bot, update *models.Update) {
	// Регистрируем или получаем пользователя
	user, err := b.getOrCreateUser(ctx, update.Message.From)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error getting or creating user")
		b.sendError(ctx, update.Message.Chat.ID, "Произошла ошибка при регистрации. Попробуйте позже.")
		return
	}

	// Приветственное сообщение
	var text string
	if user.IsAdmin {
		text = "👋 Добро пожаловать в админ-панель!\n\n🔗 Отправьте ссылку с jmmoto.ru, чтобы добавить новый мотоцикл в каталог"
	} else {
		text = "🏍️ Добро пожаловать в каталог мотоциклов!\n\n📱 Нажмите кнопку \"Каталог\" чтобы посмотреть доступные мотоциклы"
	}

	b.sendMessage(ctx, update.Message.Chat.ID, text)
}

func (b *Bot) handleMessage(ctx context.Context, _ *bot.Bot, update *models.Update) {
	// Регистрируем или получаем пользователя при любом сообщении
	user, err := b.getOrCreateUser(ctx, update.Message.From)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error getting or creating user")
		// Не отправляем ошибку пользователю, чтобы не спамить
		return
	}

	text := update.Message.Text
	if text == "" {
		return
	}

	// Проверяем, ожидаем ли мы ввод цены
	if motorcycleID, ok := b.waitingPrice.Load(update.Message.From.ID); ok {
		b.handlePriceInput(ctx, update, motorcycleID.(string))
		return
	}
	
	// Проверяем, ожидаем ли мы ввод даты прибытия
	if motorcycleID, ok := b.waitingArrivalDate.Load(update.Message.From.ID); ok {
		b.handleArrivalDateInput(ctx, update, motorcycleID.(string))
		return
	}

	// Проверяем, является ли сообщение URL (только для админов)
	if user.IsAdmin && b.isURL(text) {
		// Проверяем, что это ссылка с jmmoto.ru
		if b.isJMMotoURL(text) {
			b.handleURL(ctx, update, text)
			return
		} else {
			b.sendMessage(ctx, update.Message.Chat.ID, "⚠️ Пожалуйста, отправьте ссылку с сайта jmmoto.ru")
			return
		}
	}

	// Для любого другого сообщения показываем соответствующую подсказку
	if user.IsAdmin {
		b.sendMessage(ctx, update.Message.Chat.ID, "🔗 Отправьте ссылку с jmmoto.ru для добавления мотоцикла")
	} else {
		b.sendMessage(ctx, update.Message.Chat.ID, "📱 Нажмите кнопку \"Каталог\" чтобы посмотреть доступные мотоциклы")
	}
}

func (b *Bot) isURL(text string) bool {
	u, err := url.Parse(text)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func (b *Bot) isJMMotoURL(text string) bool {
	u, err := url.Parse(text)
	if err != nil {
		return false
	}
	return u.Host == "jmmoto.ru" || u.Host == "www.jmmoto.ru"
}

func (b *Bot) handleURL(ctx context.Context, update *models.Update, urlText string) {
	// Получаем или создаем пользователя
	user, err := b.getOrCreateUser(ctx, update.Message.From)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error getting or creating user")
		b.sendError(ctx, update.Message.Chat.ID, "Произошла ошибка при получении информации о вас. Попробуйте позже.")
		return
	}

	// Проверяем, что пользователь - админ
	if !user.IsAdmin {
		b.sendMessage(ctx, update.Message.Chat.ID, "🚫 У вас нет прав для добавления мотоциклов.\n\n📱 Используйте кнопку \"Каталог\" для просмотра доступных мотоциклов")
		return
	}

	// Отправляем сообщение о начале обработки
	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "🔄 Обрабатываю страницу и загружаю фотографии...",
	})
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error sending message")
		return
	}

	// Создаем мотоцикл из URL
	motorcycle, err := b.cases.Motorcycle.CreateMotorcycleFromURL(usecase.NewContext(ctx, user), urlText)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error creating motorcycle from URL")
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    update.Message.Chat.ID,
			MessageID: msg.ID,
			Text:       fmt.Sprintf("❌ Ошибка при обработке страницы: %v", err),
		})
		return
	}

	// Сохраняем состояние ожидания ввода цены
	b.waitingPrice.Store(update.Message.From.ID, motorcycle.ID)

	// Обновляем сообщение с просьбой ввести цену
	b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: msg.ID,
		Text: fmt.Sprintf("✅ Мотоцикл успешно добавлен:\n🏍️ %s\n\n💰 Введите цену в рублях (только число, например: 500000)",
			motorcycle.Title),
	})
}

func (b *Bot) handlePriceInput(ctx context.Context, update *models.Update, motorcycleID string) {
	// Удаляем состояние ожидания
	b.waitingPrice.Delete(update.Message.From.ID)

	// Парсим цену
	priceText := update.Message.Text
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil {
		b.sendMessage(ctx, update.Message.Chat.ID, "❌ Неверный формат цены.\n💰 Введите число, например: 500000")
		// Восстанавливаем состояние ожидания
		b.waitingPrice.Store(update.Message.From.ID, motorcycleID)
		return
	}

	if price <= 0 {
		b.sendMessage(ctx, update.Message.Chat.ID, "❌ Цена должна быть больше нуля.\n💰 Введите корректную сумму")
		b.waitingPrice.Store(update.Message.From.ID, motorcycleID)
		return
	}

	// Получаем или создаем пользователя
	user, err := b.getOrCreateUser(ctx, update.Message.From)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error getting or creating user")
		b.sendError(ctx, update.Message.Chat.ID, "Произошла ошибка при получении информации о вас.")
		return
	}

	// Обновляем только цену пока
	patch := &domain.PatchMotorcycle{
		Price: &price,
	}

	_, err = b.cases.Motorcycle.PatchMotorcycle(usecase.NewContext(ctx, user), motorcycleID, patch)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error updating motorcycle")
		b.sendError(ctx, update.Message.Chat.ID, "Ошибка при обновлении мотоцикла.")
		return
	}

	// Переходим к запросу даты прибытия
	b.waitingArrivalDate.Store(update.Message.From.ID, motorcycleID)
	b.sendMessage(ctx, update.Message.Chat.ID, "✅ Цена установлена!\n\n📅 Когда прибудет мотоцикл? (введите дату в любом формате, например: \"15 февраля\" или \"через неделю\")")
}

func (b *Bot) sendMessage(ctx context.Context, chatID int64, text string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error sending message")
	}
}

func (b *Bot) sendError(ctx context.Context, chatID int64, text string) {
	b.sendMessage(ctx, chatID, fmt.Sprintf("❌ %s", text))
}

// getOrCreateUser получает пользователя из БД или создает нового из данных Telegram
func (b *Bot) handleArrivalDateInput(ctx context.Context, update *models.Update, motorcycleID string) {
	// Удаляем состояние ожидания
	b.waitingArrivalDate.Delete(update.Message.From.ID)

	// Получаем введенную дату (без дополнительной обработки, кроме базовой проверки)
	arrivalDateText := update.Message.Text
	
	// Базовая проверка на SQL инъекцию - удаляем потенциально опасные символы
	if len(arrivalDateText) > 200 || arrivalDateText == "" {
		b.sendMessage(ctx, update.Message.Chat.ID, "❌ Дата слишком длинная или пустая.\n📅 Введите дату прибытия (например: \"15 февраля\" или \"через неделю\")")
		b.waitingArrivalDate.Store(update.Message.From.ID, motorcycleID)
		return
	}

	// Получаем или создаем пользователя
	user, err := b.getOrCreateUser(ctx, update.Message.From)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error getting or creating user")
		b.sendError(ctx, update.Message.Chat.ID, "Произошла ошибка при получении информации о вас.")
		return
	}

	// Получаем текущие данные мотоцикла
	motorcycle, err := b.cases.Motorcycle.GetMotorcycle(usecase.NewContext(ctx, user), motorcycleID)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error getting motorcycle")
		b.sendError(ctx, update.Message.Chat.ID, "Ошибка при получении данных мотоцикла.")
		return
	}

	// Обновляем данные мотоцикла с датой прибытия и устанавливаем статус available
	var motorcycleData *domain.MotorcycleData
	if motorcycle.Data != nil {
		motorcycleData = motorcycle.Data
	} else {
		motorcycleData = &domain.MotorcycleData{}
	}
	motorcycleData.ArrivalDate = arrivalDateText

	status := domain.MotorcycleStatusAvailable
	patch := &domain.PatchMotorcycle{
		Data:   motorcycleData,
		Status: &status,
	}

	updatedMotorcycle, err := b.cases.Motorcycle.PatchMotorcycle(usecase.NewContext(ctx, user), motorcycleID, patch)
	if err != nil {
		slogx.FromCtxWithErr(ctx, err).Error("error updating motorcycle")
		b.sendError(ctx, update.Message.Chat.ID, "Ошибка при обновлении мотоцикла.")
		return
	}

	b.sendMessage(ctx, update.Message.Chat.ID, fmt.Sprintf(
		"🎉 Мотоцикл успешно добавлен в каталог!\n\n🏍️ %s\n💰 Цена: %.0f ₽\n📅 Дата прибытия: %s\n📊 Статус: %s\n\n✨ Теперь он доступен в мини-приложении!",
		updatedMotorcycle.Title,
		updatedMotorcycle.Price,
		arrivalDateText,
		updatedMotorcycle.Status,
	))
}

func (b *Bot) getOrCreateUser(ctx context.Context, from *models.User) (*domain.User, error) {
	// Пытаемся получить пользователя
	user, err := b.cases.User.GetByTelegramID(ctx, from.ID)
	if err == nil {
		return user, nil
	}

	// Если пользователь не найден - создаем нового
	if errors.Is(err, repo.ErrNotFound) {
		// Формируем данные пользователя из Telegram
		firstName := from.FirstName
		lastName := ""
		username := ""
		avatar := ""

		if from.LastName != "" {
			lastName = from.LastName
		}
		if from.Username != "" {
			username = from.Username
		}
		// Avatar можно получить через GetUserProfilePhotos, но это сложнее
		// Пока оставляем пустым, будет загружен при первом использовании через Mini App

		tgData := &domain.UserTGData{
			TelegramID:       from.ID,
			FirstName:        firstName,
			LastName:         lastName,
			TelegramUsername: username,
			Avatar:           avatar,
		}

		// Создаем пользователя
		user, err = b.cases.User.Create(ctx, tgData)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		b.log.Info("new user registered", "telegram_id", from.ID, "username", username)
		return user, nil
	}

	// Другая ошибка
	return nil, fmt.Errorf("failed to get user: %w", err)
}
