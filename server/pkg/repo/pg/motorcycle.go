package pg

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shampsdev/go-telegram-template/pkg/domain"
)

type MotorcycleRepo struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewMotorcycleRepo(db *pgxpool.Pool) *MotorcycleRepo {
	return &MotorcycleRepo{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *MotorcycleRepo) Create(ctx context.Context, motorcycle *domain.CreateMotorcycle) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Создаем мотоцикл
	s := r.psql.Insert(`"motorcycle"`).
		Columns("title", "price", "currency", "description", "status", "source_url").
		Values(motorcycle.Title, motorcycle.Price, motorcycle.Currency, motorcycle.Description, motorcycle.Status, motorcycle.SourceURL).
		Suffix("RETURNING id")

	sql, args, err := s.ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to build SQL: %w", err)
	}

	var id string
	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create motorcycle: %w", err)
	}

	// Создаем фотографии
	for i, photoURL := range motorcycle.PhotoURLs {
		photoSQL := r.psql.Insert(`"motorcycle_photo"`).
			Columns("motorcycle_id", "s3_url", "\"order\"").
			Values(id, photoURL, i).
			Suffix("RETURNING id")

		photoSQLStr, photoArgs, err := photoSQL.ToSql()
		if err != nil {
			return "", fmt.Errorf("failed to build photo SQL: %w", err)
		}

		var photoID string
		err = tx.QueryRow(ctx, photoSQLStr, photoArgs...).Scan(&photoID)
		if err != nil {
			return "", fmt.Errorf("failed to create photo: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (r *MotorcycleRepo) Patch(ctx context.Context, id string, motorcycle *domain.PatchMotorcycle) error {
	s := r.psql.Update(`"motorcycle"`).
		Where(sq.Eq{"id": id}).
		Set("updated_at", time.Now())

	if motorcycle.Title != nil {
		s = s.Set("title", *motorcycle.Title)
	}
	if motorcycle.Price != nil {
		s = s.Set("price", *motorcycle.Price)
	}
	if motorcycle.Currency != nil {
		s = s.Set("currency", *motorcycle.Currency)
	}
	if motorcycle.Description != nil {
		s = s.Set("description", *motorcycle.Description)
	}
	if motorcycle.Status != nil {
		s = s.Set("status", *motorcycle.Status)
	}

	sql, args, err := s.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *MotorcycleRepo) Filter(ctx context.Context, filter *domain.FilterMotorcycle) ([]*domain.Motorcycle, error) {
	s := r.psql.Select("m.id", "m.title", "m.price", "m.currency", "m.description", "m.status", "m.source_url", "m.created_at", "m.updated_at").
		From(`"motorcycle" m`)

	if filter.ID != nil {
		s = s.Where(sq.Eq{"m.id": *filter.ID})
	}
	if filter.Status != nil {
		s = s.Where(sq.Eq{"m.status": *filter.Status})
	}
	if filter.Title != nil {
		s = s.Where(sq.Like{"m.title": "%" + *filter.Title + "%"})
	}
	if filter.MinPrice != nil {
		s = s.Where(sq.GtOrEq{"m.price": *filter.MinPrice})
	}
	if filter.MaxPrice != nil {
		s = s.Where(sq.LtOrEq{"m.price": *filter.MaxPrice})
	}

	// Сортировка: available -> reserved -> sold
	s = s.OrderBy(`
		CASE 
			WHEN m.status = 'available' THEN 1
			WHEN m.status = 'reserved' THEN 2
			WHEN m.status = 'sold' THEN 3
			ELSE 4
		END,
		m.created_at DESC
	`)

	sql, args, err := s.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL: %w", err)
	}
	defer rows.Close()

	motorcycles := []*domain.Motorcycle{}
	motorcycleMap := make(map[string]*domain.Motorcycle)

	for rows.Next() {
		var m domain.Motorcycle
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Price,
			&m.Currency,
			&m.Description,
			&m.Status,
			&m.SourceURL,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		motorcycles = append(motorcycles, &m)
		motorcycleMap[m.ID] = &m
	}

	// Загружаем фотографии, если нужно
	if filter.IncludePhotos && len(motorcycles) > 0 {
		motorcycleIDs := make([]string, 0, len(motorcycles))
		for _, m := range motorcycles {
			motorcycleIDs = append(motorcycleIDs, m.ID)
		}

		photoSQL := r.psql.Select("id", "motorcycle_id", "s3_url", "\"order\"", "created_at").
			From(`"motorcycle_photo"`).
			Where(sq.Eq{"motorcycle_id": motorcycleIDs}).
			OrderBy("\"order\" ASC")

		photoSQLStr, photoArgs, err := photoSQL.ToSql()
		if err != nil {
			return nil, fmt.Errorf("failed to build photo SQL: %w", err)
		}

		photoRows, err := r.db.Query(ctx, photoSQLStr, photoArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to query photos: %w", err)
		}
		defer photoRows.Close()

		for photoRows.Next() {
			var photo domain.MotorcyclePhoto
			err := photoRows.Scan(
				&photo.ID,
				&photo.MotorcycleID,
				&photo.S3URL,
				&photo.Order,
				&photo.CreatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan photo row: %w", err)
			}

			if m, ok := motorcycleMap[photo.MotorcycleID]; ok {
				if m.Photos == nil {
					m.Photos = []*domain.MotorcyclePhoto{}
				}
				m.Photos = append(m.Photos, &photo)
			}
		}
	}

	return motorcycles, nil
}

func (r *MotorcycleRepo) Delete(ctx context.Context, id string) error {
	s := r.psql.Delete(`"motorcycle"`).
		Where(sq.Eq{"id": id})

	sql, args, err := s.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *MotorcycleRepo) AddPhotos(ctx context.Context, motorcycleID string, photoURLs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Получаем текущий максимальный порядок фотографий
	maxOrderSQL := r.psql.Select("COALESCE(MAX(\"order\"), -1)").
		From(`"motorcycle_photo"`).
		Where(sq.Eq{"motorcycle_id": motorcycleID})

	maxOrderSQLStr, maxOrderArgs, err := maxOrderSQL.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build max order SQL: %w", err)
	}

	var maxOrder int
	err = tx.QueryRow(ctx, maxOrderSQLStr, maxOrderArgs...).Scan(&maxOrder)
	if err != nil {
		return fmt.Errorf("failed to get max order: %w", err)
	}

	// Добавляем новые фотографии
	for i, photoURL := range photoURLs {
		order := maxOrder + 1 + i
		photoSQL := r.psql.Insert(`"motorcycle_photo"`).
			Columns("motorcycle_id", "s3_url", "\"order\"").
			Values(motorcycleID, photoURL, order).
			Suffix("RETURNING id")

		photoSQLStr, photoArgs, err := photoSQL.ToSql()
		if err != nil {
			return fmt.Errorf("failed to build photo SQL: %w", err)
		}

		var photoID string
		err = tx.QueryRow(ctx, photoSQLStr, photoArgs...).Scan(&photoID)
		if err != nil {
			return fmt.Errorf("failed to create photo: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

