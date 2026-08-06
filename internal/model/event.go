package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrEmptyEventID    = errors.New("event_id field is required")
	ErrInvalidUUID     = errors.New("event_id must be a valid UUID")
	ErrInvalidUserID   = errors.New("user_id must be greater than zero")
	ErrEmptyAction     = errors.New("action field is required")
	ErrZeroTimestamp   = errors.New("timestamp must not be zero")
	ErrFutureTimestamp = errors.New("timestamp cannot be in the future")
)

type ActivityEvent struct {
	EventID          string         `json:"event_id"`
	UserID           int64          `json:"user_id"`
	Action           string         `json:"action"`
	ActionObjectID   int64          `json:"action_object_id,omitempty"`
	ActionObjectType string         `json:"action_object_type,omitempty"`
	Timestamp        time.Time      `json:"timestamp"`
	Metadata         map[string]any `json:"metadata,omitempty"`
}

func (event *ActivityEvent) Validate() error {
	const fn = "validator.ValidateEventStruct"

	eventID := strings.TrimSpace(event.EventID)
	if eventID == "" {
		return fmt.Errorf("%s: %w", fn, ErrEmptyEventID)
	}

	if _, err := uuid.Parse(eventID); err != nil {
		return fmt.Errorf("%s: %w (%v)", fn, ErrInvalidUUID, err)
	}

	if event.UserID <= 0 {
		return fmt.Errorf("%s: %w", fn, ErrInvalidUserID)
	}

	action := strings.TrimSpace(event.Action)
	if action == "" {
		return fmt.Errorf("%s: %w", fn, ErrEmptyAction)
	}

	if event.Timestamp.IsZero() {
		return fmt.Errorf("%s: %w", fn, ErrZeroTimestamp)
	}

	if event.Timestamp.After(time.Now().Add(5 * time.Second)) {
		return fmt.Errorf("%s: %w", fn, ErrFutureTimestamp)
	}

	if event.ActionObjectType != "" && event.ActionObjectID <= 0 {
		return fmt.Errorf("%s: action_object_id must be greater than zero when action_object_type is provided", fn)
	}

	if event.ActionObjectID > 0 && strings.TrimSpace(event.ActionObjectType) == "" {
		return fmt.Errorf("%s: action_object_type must be provided when action_object_id is specified", fn)
	}

	return nil
}
