package pg

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shampsdev/go-telegram-template/pkg/domain"
)

type AnalyticsRepo struct {
	db   *pgxpool.Pool
	psql squirrel.StatementBuilderType
}

func NewAnalyticsRepo(db *pgxpool.Pool) *AnalyticsRepo {
	return &AnalyticsRepo{
		db:   db,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *AnalyticsRepo) RecordVisit(ctx context.Context, visit *domain.CreateUserVisit) error {
	s := r.psql.Insert("user_visits").
		Columns("user_id", "session_id", "source").
		Values(visit.UserID, visit.SessionID, visit.Source)

	sql, args, err := s.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to record visit: %w", err)
	}

	return nil
}

func (r *AnalyticsRepo) GetUserStats(ctx context.Context, userID string) (*domain.UserVisitStats, error) {
	s := r.psql.Select("*").
		From("user_visit_stats").
		Where(squirrel.Eq{"user_id": userID})

	sql, args, err := s.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	var stats domain.UserVisitStats
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&stats.UserID,
		&stats.TotalVisits,
		&stats.UniqueDays,
		&stats.FirstVisit,
		&stats.LastVisit,
		&stats.DaysSpan,
		&stats.AvgVisitsPerDay,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return &stats, nil
}

func (r *AnalyticsRepo) GetAllUserStats(ctx context.Context) ([]*domain.UserVisitStats, error) {
	s := r.psql.Select("*").
		From("user_visit_stats").
		OrderBy("total_visits DESC")

	sql, args, err := s.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query user stats: %w", err)
	}
	defer rows.Close()

	var stats []*domain.UserVisitStats
	for rows.Next() {
		var stat domain.UserVisitStats
		err := rows.Scan(
			&stat.UserID,
			&stat.TotalVisits,
			&stat.UniqueDays,
			&stat.FirstVisit,
			&stat.LastVisit,
			&stat.DaysSpan,
			&stat.AvgVisitsPerDay,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user stat: %w", err)
		}
		stats = append(stats, &stat)
	}

	return stats, nil
}
