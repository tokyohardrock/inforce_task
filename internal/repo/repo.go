package repo

import (
	"context"
	"inforce_task/internal/model"
	"time"
)

type Repo interface {
	Save(ctx context.Context, event model.ActivityEvent) error
	GetByUserID(ctx context.Context, userID int64, from, to time.Time) ([]model.ActivityEvent, error)
}
