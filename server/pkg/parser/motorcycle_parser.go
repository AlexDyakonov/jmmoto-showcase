package parser

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/shampsdev/go-telegram-template/pkg/domain"
)

type MotorcycleParser struct{}

func NewMotorcycleParser() *MotorcycleParser {
	return &MotorcycleParser{}
}

func (p *MotorcycleParser) ParseMotorcycle(url string) (*domain.ParsedMotorcycleData, error) {
	// Создаем HTTP клиент
	client := &http.Client{}

	// Создаем запрос
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Устанавливаем User-Agent для имитации браузера
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неверный статус код: %d", resp.StatusCode)
	}

	// Парсим HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга HTML: %w", err)
	}

	data := &domain.ParsedMotorcycleData{}

	// Получаем весь текст страницы один раз для парсинга
	bodyText := doc.Find("body").Text()
	bodyHTML, _ := doc.Find("body").Html()

	// Парсим название мотоцикла
	nameSelectors := []string{
		"div[class*='styles_text'][class*='styles_weight--semi-bold']",
		"div[class*='styles_text__'][class*='styles_weight--semi-bold']",
		"div[class*='styles_text'][class*='uppercase']",
	}

	for _, selector := range nameSelectors {
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 5 && len(text) < 100 {
				brands := []string{"suzuki", "yamaha", "honda", "kawasaki", "ducati", "bmw", "triumph", "harley"}
				textLower := strings.ToLower(text)
				for _, brand := range brands {
					if strings.Contains(textLower, brand) {
						text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
						data.Name = strings.TrimSpace(text)
						return
					}
				}
			}
		})
		if data.Name != "" {
			break
		}
	}

	// Если не нашли в специальном div, ищем в заголовках
	if data.Name == "" {
		nameSelector := "h1, h2"
		doc.Find(nameSelector).Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" && len(text) > 5 && len(text) < 100 {
				brands := []string{"suzuki", "yamaha", "honda", "kawasaki", "ducati", "bmw", "triumph", "harley"}
				textLower := strings.ToLower(text)
				for _, brand := range brands {
					if strings.Contains(textLower, brand) {
						data.Name = text
						return
					}
				}
			}
		})
	}

	// Если не нашли, ищем в HTML с помощью regex
	if data.Name == "" {
		nameRe := regexp.MustCompile(`(Suzuki|Yamaha|Honda|Kawasaki|Ducati|BMW|Triumph|Harley-Davidson|Harley)(?:<!--\s*-->)?\s*(?:<!--\s*-->)?\s*([A-Z0-9]+)`)
		matches := nameRe.FindStringSubmatch(bodyHTML)
		if len(matches) >= 3 {
			data.Name = strings.TrimSpace(matches[1] + " " + matches[2])
		}
	}

	// Парсим год
	yearRe := regexp.MustCompile(`Год:\s*(\d{4})`)
	if matches := yearRe.FindStringSubmatch(bodyText); len(matches) >= 2 {
		if year, err := strconv.Atoi(matches[1]); err == nil {
			data.Year = year
		}
	}

	// Парсим пробег
	mileageRe := regexp.MustCompile(`Пробег:\s*(\d+)\s*км`)
	if matches := mileageRe.FindStringSubmatch(bodyText); len(matches) >= 2 {
		if mileage, err := strconv.Atoi(matches[1]); err == nil {
			data.Mileage = mileage
		}
	}

	// Парсим объем
	volumeRe := regexp.MustCompile(`Объем:\s*(\d+)\s*сс`)
	if matches := volumeRe.FindStringSubmatch(bodyText); len(matches) >= 2 {
		if volume, err := strconv.Atoi(matches[1]); err == nil {
			data.Volume = volume
		}
	}

	// Парсим номер рамы
	frameRe := regexp.MustCompile(`Номер рамы:\s*([A-Z0-9-]+)`)
	if matches := frameRe.FindStringSubmatch(bodyText); len(matches) >= 2 {
		data.FrameNum = strings.TrimSpace(matches[1])
	}

	// Парсим изображения
	imageMap := make(map[string]bool)

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		attrs := []string{"src", "data-src", "data-lazy-src"}

		for _, attr := range attrs {
			src, exists := s.Attr(attr)
			if !exists || src == "" {
				continue
			}

			// Пропускаем логотипы, иконки и другие служебные изображения
			srcLower := strings.ToLower(src)
			skipPatterns := []string{"logo", "icon", "vite.svg", "favicon", "sprite", "placeholder", "search-banner"}
			shouldSkip := false
			for _, pattern := range skipPatterns {
				if strings.Contains(srcLower, pattern) {
					shouldSkip = true
					break
				}
			}
			if shouldSkip {
				continue
			}

			// Преобразуем относительные URL в абсолютные
			if strings.HasPrefix(src, "//") {
				src = "https:" + src
			} else if strings.HasPrefix(src, "/") {
				src = "https://jmmoto.ru" + src
			} else if !strings.HasPrefix(src, "http") {
				src = "https://jmmoto.ru/" + src
			}

			// Добавляем только уникальные изображения
			if !imageMap[src] {
				imageMap[src] = true
				data.Images = append(data.Images, src)
			}
		}
	})

	return data, nil
}

