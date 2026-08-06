package handler

import (
	"encoding/json"
	"fmt"
	"inforce_task/internal/model"
	"log/slog"
	"net/http"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	event := new(model.ActivityEvent)

	err := json.NewDecoder(r.Body).Decode(event)
	if err != nil {
		slog.Error(err.Error())

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = event.Validate()
	if err != nil {
		slog.Error(err.Error())

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(*event)

	w.WriteHeader(http.StatusOK)
}
