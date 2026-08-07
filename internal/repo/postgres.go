package repo

import (
	"context"
	"fmt"
	"time"

	"inforce_task/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresEventRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresEventRepository(pool *pgxpool.Pool) *PostgresEventRepository {
	return &PostgresEventRepository{pool: pool}
}

func (r *PostgresEventRepository) Save(ctx context.Context, event model.ActivityEvent) error {
	const fn = "repository.PostgresEventRepository.Save"

	query := `
		INSERT INTO events (
			event_id,
			user_id,
			action,
			action_object_id,
			action_object_type,
			timestamp,
			metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (event_id) DO NOTHING`

	_, err := r.pool.Exec(
		ctx,
		query,
		event.EventID,
		event.UserID,
		event.Action,
		event.ActionObjectID,
		event.ActionObjectType,
		event.Timestamp,
		event.Metadata,
	)
	if err != nil {
		return fmt.Errorf("%s: exec insert: %w", fn, err)
	}

	return nil
}

func (r *PostgresEventRepository) GetByUserID(ctx context.Context, userID int64, from, to time.Time) ([]model.ActivityEvent, error) {
	const fn = "repository.PostgresEventRepository.GetByUserID"

	query := `
		SELECT
			event_id,
			user_id,
			action,
			COALESCE(action_object_id, 0) AS action_object_id,
			COALESCE(action_object_type, '') AS action_object_type,
			timestamp,
			metadata
		FROM events
		WHERE user_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp DESC`

	rows, err := r.pool.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("%s: query events: %w", fn, err)
	}

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ActivityEvent])
	if err != nil {
		return nil, fmt.Errorf("%s: collect rows: %w", fn, err)
	}

	return events, nil
}

func (r *PostgresEventRepository) CalculateAndSaveUserStats(ctx context.Context, from, to time.Time) error {
	const fn = "repo.PostgresEventRepository.CalculateAndSaveUserStats"

	query := `
		INSERT INTO user_activity_stats (user_id, time_bucket, event_count)
		SELECT
			user_id,
			$1::timestamptz AS time_bucket,
			COUNT(*) AS event_count
		FROM events
		WHERE timestamp >= $1 AND timestamp < $2
		GROUP BY user_id
		ON CONFLICT (user_id, time_bucket)
		DO UPDATE SET event_count = EXCLUDED.event_count`

	_, err := r.pool.Exec(ctx, query, from, to)
	if err != nil {
		return fmt.Errorf("%s: exec insert stats: %w", fn, err)
	}

	return nil
}

func (r *PostgresEventRepository) GetUserStats(ctx context.Context, userID int64, from, to time.Time) (*model.UserStats, error) {
	const fn = "repo.PostgresEventRepository.GetUserStats"

	query := `
		SELECT
			COALESCE(SUM(event_count), 0)::int AS total_events
		FROM user_activity_stats
		WHERE user_id = $1
		  AND time_bucket >= date_bin('4 hours', $2::timestamptz, '2000-01-01 00:00:00Z')
		  AND time_bucket < $3`

	stats := model.UserStats{
		UserID:    userID,
		StartTime: from,
		EndTime:   to,
	}

	err := r.pool.QueryRow(ctx, query, userID, from, to).Scan(
		&stats.TotalEvents,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: exec select stats: %w", fn, err)
	}

	return &stats, nil
}
