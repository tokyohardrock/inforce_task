package handler

import (
	"encoding/json"
	"inforce_task/internal/model"
	"inforce_task/internal/repo"
	"log/slog"
	"mime"
	"net/http"
	"strconv"
	"time"
)

type EventHandler struct {
	repo repo.Repo
}

func NewEventHandler(repo repo.Repo) *EventHandler {
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
		slog.Error(err.Error())

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = event.Validate()
	if err != nil {
		slog.Error(err.Error())

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
