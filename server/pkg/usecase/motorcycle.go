package usecase

import (
	"context"
	"fmt"

	"github.com/shampsdev/go-telegram-template/pkg/domain"
	"github.com/shampsdev/go-telegram-template/pkg/repo"
)

type Motorcycle struct {
	motorcycleRepo repo.Motorcycle
	storage        repo.ImageStorage
	parser         MotorcycleParser
}

type MotorcycleParser interface {
	ParseMotorcycle(url string) (*domain.ParsedMotorcycleData, error)
}

func NewMotorcycle(motorcycleRepo repo.Motorcycle, storage repo.ImageStorage, parser MotorcycleParser) *Motorcycle {
	return &Motorcycle{
		motorcycleRepo: motorcycleRepo,
		storage:        storage,
		parser:         parser,
	}
}

func (m *Motorcycle) ListMotorcycles(ctx context.Context, filter *domain.FilterMotorcycle) ([]*domain.Motorcycle, error) {
	if filter == nil {
		filter = &domain.FilterMotorcycle{}
	}
	filter.IncludePhotos = true
	return m.motorcycleRepo.Filter(ctx, filter)
}

func (m *Motorcycle) GetMotorcycle(ctx context.Context, id string) (*domain.Motorcycle, error) {
	filter := &domain.FilterMotorcycle{
		ID:            &id,
		IncludePhotos: true,
	}
	return repo.First(m.motorcycleRepo.Filter)(ctx, filter)
}

func (m *Motorcycle) CreateMotorcycle(ctx Context, createMotorcycle *domain.CreateMotorcycle) (*domain.Motorcycle, error) {
	// Сохраняем исходные URL фотографий
	originalPhotoURLs := createMotorcycle.PhotoURLs
	
	// Сначала создаем мотоцикл с пустым массивом фотографий, чтобы получить ID
	createMotorcycle.PhotoURLs = []string{}
	id, err := m.motorcycleRepo.Create(ctx, createMotorcycle)
	if err != nil {
		return nil, fmt.Errorf("failed to create motorcycle: %w", err)
	}

	// Теперь загружаем фотографии в S3 с правильным ключом на основе ID мотоцикла
	photoURLs := make([]string, 0, len(originalPhotoURLs))
	for i, photoURL := range originalPhotoURLs {
		// Используем ID мотоцикла и индекс фотографии для создания уникального ключа
		key := fmt.Sprintf("motorcycles/%s/%d", id, i)
		s3URL, err := m.storage.SaveImageByURL(ctx, photoURL, key)
		if err != nil {
			return nil, fmt.Errorf("failed to save photo %d: %w", i, err)
		}
		photoURLs = append(photoURLs, s3URL)
	}

	// Добавляем фотографии в БД
	if len(photoURLs) > 0 {
		err = m.motorcycleRepo.AddPhotos(ctx, id, photoURLs)
		if err != nil {
			return nil, fmt.Errorf("failed to add photos: %w", err)
		}
	}

	return m.GetMotorcycle(ctx, id)
}

func (m *Motorcycle) CreateMotorcycleFromURL(ctx Context, url string) (*domain.Motorcycle, error) {
	// Парсим страницу
	data, err := m.parser.ParseMotorcycle(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse motorcycle page: %w", err)
	}

	// Формируем название из данных парсера
	title := data.Name
	if data.Year > 0 {
		title = fmt.Sprintf("%s %d", title, data.Year)
	}

	// Создаем мотоцикл со статусом draft
	createMotorcycle := &domain.CreateMotorcycle{
		Title:     title,
		Price:     0, // Цену нужно будет ввести позже
		Currency:  "RUB",
		Status:    domain.MotorcycleStatusDraft,
		SourceURL: url,
		PhotoURLs: data.Images,
	}

	// Добавляем структурированные данные, если есть
	if data.Mileage > 0 || data.Volume > 0 || data.FrameNum != "" {
		motorcycleData := &domain.MotorcycleData{}
		
		if data.Mileage > 0 {
			motorcycleData.Mileage = &data.Mileage
			motorcycleData.MileageUnit = "км"
		}
		if data.Volume > 0 {
			motorcycleData.Volume = &data.Volume
			motorcycleData.VolumeUnit = "сс"
		}
		if data.FrameNum != "" {
			motorcycleData.FrameNumber = data.FrameNum
		}
		
		createMotorcycle.Data = motorcycleData
	}

	return m.CreateMotorcycle(ctx, createMotorcycle)
}

func (m *Motorcycle) PatchMotorcycle(ctx Context, id string, patchMotorcycle *domain.PatchMotorcycle) (*domain.Motorcycle, error) {
	err := m.motorcycleRepo.Patch(ctx, id, patchMotorcycle)
	if err != nil {
		return nil, fmt.Errorf("failed to patch motorcycle: %w", err)
	}
	return m.GetMotorcycle(ctx, id)
}

func (m *Motorcycle) UpdateMotorcycleStatus(ctx Context, id string, status domain.MotorcycleStatus) (*domain.Motorcycle, error) {
	patch := &domain.PatchMotorcycle{
		Status: &status,
	}
	return m.PatchMotorcycle(ctx, id, patch)
}

