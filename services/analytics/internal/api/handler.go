package api

import (
	"encoding/json"
	"net/http"

	chclient "github.com/palach608/go-platform/services/analytics/internal/clickhouse"
)

type Handler struct {
	ch *chclient.Client
}

func NewHandler(ch *chclient.Client) *Handler {
	return &Handler{ch: ch}
}

func (h *Handler) TopUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.ch.TopUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) NotesPerDay(w http.ResponseWriter, r *http.Request) {
	days, err := h.ch.NotesPerDay(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(days)
}
