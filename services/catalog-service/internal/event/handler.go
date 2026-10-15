package event

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	events, err := h.repo.GetEvents(r.Context())
	if err != nil {
		slog.Error("failed to get events", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to get events")
		return
	}

	writeJSON(w, http.StatusOK, events)
}

func (h *Handler) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		slog.Error("id is required")
		writeJSONError(w, http.StatusBadRequest, "id is required")
		return
	}

	event, err := h.repo.GetEvent(r.Context(), id)
	if err != nil {
		slog.Error("failed to get event", "error", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to get event")
	}
	writeJSON(w, http.StatusOK, event)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
