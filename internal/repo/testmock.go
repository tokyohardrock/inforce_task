package repo

import (
	"context"
	"inforce_task/internal/model"
	"time"
)

type MockDB struct {
	Err error
}

func NewMockEventRepository() *MockDB {
	return &MockDB{}
}

func (m *MockDB) Save(ctx context.Context, event model.ActivityEvent) error {
	return m.Err
}

func (m *MockDB) GetByUserID(ctx context.Context, userID int64, from, to time.Time) ([]model.ActivityEvent, error) {
	return nil, m.Err
}
