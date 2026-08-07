package handler

import (
	"context"
	"encoding/json"
	"inforce_task/internal/model"
	"log/slog"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Repo interface {
	Save(ctx context.Context, event model.ActivityEvent) error
	GetByUserID(ctx context.Context, userID int64, from, to time.Time) ([]model.ActivityEvent, error)
	GetUserStats(ctx context.Context, userID int64, from, to time.Time) (*model.UserStats, error)
}

type EventHandler struct {
	repo Repo
}

func NewEventHandler(repo Repo) *EventHandler {
	return &EventHandler{
		repo: repo,
	}
}

func IsJSONContentType(contentType string) bool {
	if contentType == "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	return mediaType == "application/json"
}

func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if !IsJSONContentType(r.Header.Get("Content-Type")) {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	event := new(model.ActivityEvent)

	err := json.NewDecoder(r.Body).Decode(event)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = event.Validate()
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.repo.Save(r.Context(), *event)
	if err != nil {
		slog.Error(err.Error())

		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userIDStr := strings.TrimSpace(query.Get("user_id"))
	startTimeStr := strings.TrimSpace(query.Get("start_time"))
	endTimeStr := strings.TrimSpace(query.Get("end_time"))

	if userIDStr == "" {
		http.Error(w, "query param user_id is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid user_id format", http.StatusBadRequest)
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		http.Error(w, "invalid start_time format (expected RFC3339, e.g. 2026-08-06T12:00:00Z)", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		http.Error(w, "invalid end_time format (expected RFC3339, e.g. 2026-08-06T12:00:00Z)", http.StatusBadRequest)
		return
	}

	if startTime.After(endTime) {
		http.Error(w, "start_time must be before or equal to end_time", http.StatusBadRequest)
		return
	}

	events, err := h.repo.GetByUserID(r.Context(), userID, startTime, endTime)
	if err != nil {
		slog.Error("failed to get events from db", "user_id", userID, "error", err)

		http.Error(w, "unable to get events", http.StatusInternalServerError)
		return
	}

	if events == nil {
		events = []model.ActivityEvent{}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(events); err != nil {
		slog.Error("failed to encode events response", "error", err)
	}
}

func (h *EventHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		http.Error(w, "invalid or missing user_id", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	startTimeStr := strings.TrimSpace(query.Get("start_time"))
	endTimeStr := strings.TrimSpace(query.Get("end_time"))

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		http.Error(w, "invalid start_time format (expected RFC3339)", http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		http.Error(w, "invalid end_time format (expected RFC3339)", http.StatusBadRequest)
		return
	}

	if startTime.After(endTime) {
		http.Error(w, "start_time must be before or equal to end_time", http.StatusBadRequest)
		return
	}

	stats, err := h.repo.GetUserStats(r.Context(), userID, startTime, endTime)
	if err != nil {
		slog.Error("failed to fetch user stats", "user_id", userID, "error", err)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		slog.Error("failed to encode stats response", "error", err)
	}
}
