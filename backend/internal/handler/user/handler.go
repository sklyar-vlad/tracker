package user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sklyar-vlad/selfDev/internal/handler/user/dto"
	"github.com/sklyar-vlad/selfDev/internal/model"
)

type Service interface {
	Register(ctx context.Context, username, email, password string) (model.User, error)
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	user, err := h.service.Register(r.Context(), input.Username, input.Email, input.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(dto.ToUserResponse(user)); err != nil {
		log.Printf("Invalid encode JSON: %v", err)
	}
}
