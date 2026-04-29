package api

import (
	"encoding/json"
	"net/http"

	apiModel "github.com/SilverName608/go-notes/internal/api/model"
	"github.com/SilverName608/go-notes/internal/domain/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req apiModel.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Incorrect JSON", http.StatusBadRequest)
		return
	}

	response, err := uh.svc.Register(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req apiModel.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Incorrect JSON", http.StatusBadRequest)
		return
	}

	response, err := uh.svc.Login(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
