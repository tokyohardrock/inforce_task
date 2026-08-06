package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"inforce_task/internal/model"
	"inforce_task/internal/repo"
)

func TestIsJSONContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{
			name:        "Valid simple application/json",
			contentType: "application/json",
			want:        true,
		},
		{
			name:        "Valid application/json with charset parameter",
			contentType: "application/json; charset=utf-8",
			want:        true,
		},
		{
			name:        "Valid application/json with uppercase letters and spaces",
			contentType: "Application/JSON; CHARSET=UTF-8",
			want:        true,
		},
		{
			name:        "Empty content-type header",
			contentType: "",
			want:        false,
		},
		{
			name:        "Unsupported content-type text/html",
			contentType: "text/html",
			want:        false,
		},
		{
			name:        "Unsupported content-type application/x-www-form-urlencoded",
			contentType: "application/x-www-form-urlencoded",
			want:        false,
		},
		{
			name:        "Malformed content-type header",
			contentType: "application/json; invalid-param",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsJSONContentType(tt.contentType)
			if got != tt.want {
				t.Errorf("isJSONContentType(%q) = %v; want %v", tt.contentType, got, tt.want)
			}
		})
	}
}

func TestCreateEvent(t *testing.T) {
	validEvent := model.ActivityEvent{
		EventID:   "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
		UserID:    42,
		Action:    "click",
		Timestamp: time.Now(),
	}

	validJSON, _ := json.Marshal(validEvent)

	tests := []struct {
		name           string
		contentType    string
		body           []byte
		expectedStatus int
	}{
		{
			name:           "Success - Valid JSON and Content-Type",
			contentType:    "application/json",
			body:           validJSON,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Success - Valid JSON with charset in Content-Type",
			contentType:    "application/json; charset=utf-8",
			body:           validJSON,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error - Unsupported Content-Type (text/plain)",
			contentType:    "text/plain",
			body:           validJSON,
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Error - Missing Content-Type Header",
			contentType:    "",
			body:           validJSON,
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Error - Invalid JSON syntax",
			contentType:    "application/json",
			body:           []byte(`{"event_id": "123", "user_id": invalid_json}`),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - Empty Request Body",
			contentType:    "application/json",
			body:           []byte(""),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error - Failed Event Validation (e.g. empty Action/UserID <= 0)",
			contentType:    "application/json",
			body:           []byte(`{"event_id": "123", "user_id": -1, "action": ""}`),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			rr := httptest.NewRecorder()

			eventHandler := NewEventHandler(repo.NewMockEventRepository())

			eventHandler.CreateEvent(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("CreateEvent returned wrong status code: got %v want %v",
					rr.Code, tt.expectedStatus)
			}
		})
	}
}
