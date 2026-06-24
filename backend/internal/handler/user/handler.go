package user

import (
	// "github.com/sklyar-vlad/selfDev/internal/handler/user/dto"
	// "github.com/sklyar-vlad/selfDev/internal/model/user"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/sklyar-vlad/selfDev/internal/handler/user/dto"
	model "github.com/sklyar-vlad/selfDev/internal/model/user"
)

// GetUsers(w http.ResponseWriter, r *http.Request)
// CreateUser(w http.ResponseWriter, r *http.Request)
// GetUser(w http.ResponseWriter, r *http.Request)
// DeleteUser(w http.ResponseWriter, r *http.Request)
// UpdateUser(w http.ResponseWriter, r *http.Request)

type Service interface {
	CreateUser(ctx context.Context, username, email, password string) (model.User, error)
	GetByLogin(ctx context.Context, username, email string) (model.User, error)
	GetById(ctx context.Context, userId uuid.UUID) (model.User, error)
}

type handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, err := uuid.Parse(id)
	if err != nil {
		h.logger.Error("invalid id", zap.String("id", id), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	user, err := h.service.GetById(r.Context(), userId)
	if err != nil {
		h.logger.Error("failed get user by id", zap.String("id", userId.String()), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	if err = json.NewEncoder(w).Encode(dto.ToUserResponse(user)); err != nil {
		h.logger.Error("failed create response with user", zap.String("email", user.Email), zap.Error(err))
	}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input dto.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed decode request", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(r.Context(), input.Username, input.Email, input.Password)
	if err != nil {
		h.logger.Error("failed create user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(dto.ToUserResponse(user)); err != nil {
		h.logger.Error("failed create response with user", zap.String("email", input.Email), zap.Error(err))
	}
}

// func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
// 	var input dto.UserRequest

// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		h.logger.Error("failed decode request", zap.Error(err))
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 	}

// 	user, err := h.service.CreateUser(ctx, input.Username, input.Email, input.Password)
// }

// func (h *handler) (w http.ResponseWriter, r *http.Request) {
// 	var input dto.UserRequest

// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		h.logger.Error("decode request failed", zap.Error(err))
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}

// 	if input.Email == "" {
// 		h.logger.Error("invalid email address", zap.Error(errors.ErrInvalidEmail))
// 		http.Error(w, errors.ErrInvalidEmail.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if len(input.Password) < 6 {
// 		h.logger.Error("invalid password, it should be more than 6 symbols", zap.Error(errors.ErrInvalidPassword))
// 		http.Error(w, errors.ErrInvalidPassword.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	user, err := h.service.Register(r.Context(), input.Username, input.Email, input.Password)
// 	if err != nil {
// 		h.logger.Error("failed create user", zap.Error(err))
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)

// 	if err = json.NewEncoder(w).Encode(dto.ToUserResponse(user)); err != nil {
// 		h.logger.Error("failed create response with user model", zap.String("email", user.Email), zap.Error(err))
// 	}
// }

// func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
// 	var input dto.UserRequest

// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		h.logger.Error("decode request failed", zap.Error(err))
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}

// 	if input.Email == "" {
// 		h.logger.Error("invalid email address", zap.Error(errors.ErrInvalidEmail))
// 		http.Error(w, errors.ErrInvalidEmail.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	accessToken, refreshToken, err := h.service.Login(r.Context(), input.Username, input.Email, input.Password)
// 	if err != nil {
// 		h.logger.Error("failed authorization", zap.Error(err))
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		MaxAge:   12 * 60 * 60,
// 	})

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		MaxAge:   30 * 24 * 60 * 60,
// 	})
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// }

// func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
// 	var input dto.TokenRequest

// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		h.logger.Error("decode request failed", zap.Error(err))
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}

// 	h.logger.Info("access token", zap.String("access token", input.AccessToken))
// 	accessToken, err := h.service.Refresh(r.Context(), input.AccessToken, input.RefreshToken)
// 	if err != nil {
// 		h.logger.Error("failed authorization", zap.Error(err))
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "access_token",
// 		Value:    accessToken,
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		MaxAge:   12 * 60 * 60,
// 	})

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// }
