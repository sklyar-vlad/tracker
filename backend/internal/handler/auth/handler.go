package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	appErrors "github.com/sklyar-vlad/selfDev/internal/errors"
	"github.com/sklyar-vlad/selfDev/internal/handler/auth/dto"
	model "github.com/sklyar-vlad/selfDev/internal/model/auth"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, username, email, password string) (model.Tokens, error)
	ConfirmEmail(ctx context.Context, token string) error
	// Logout(ctx context.Context, refreshToken string) error
	// Refresh(ctx context.Context, accessToken, refreshToken string) (string, error)
}

type handler struct {
	service AuthService
	logger  *zap.Logger
}

func NewHandler(service AuthService, logger *zap.Logger) *handler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed decode request", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	
	if err := h.service.Register(r.Context(), input.Username, input.Email, input.Password); err != nil {
		h.logger.Error("failed create user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	if token == "" {
		h.logger.Error("invalid token", zap.String("token", token))
		http.Error(w, "invalid token", http.StatusBadRequest)
	}

	err := h.service.ConfirmEmail(r.Context(), token)

	if errors.Is(err, appErrors.ErrTokenWasExpired) {
		h.logger.Error("token was expired", zap.Error(appErrors.ErrTokenWasExpired))
		http.Error(w, appErrors.ErrTokenWasExpired.Error(), http.StatusGone)
	}

	if err != nil {
		h.logger.Error("failed verify email", zap.Error(err))
		http.Error(w, "failed verify email", http.StatusInternalServerError)
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("failed decode request", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	tokens, err := h.service.Login(r.Context(), input.Username, input.Email, input.Password)
	if err != nil {
		h.logger.Error("failed login", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   12 * 60 * 60,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   30 * 24 * 60 * 60,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
// 	var input dto.LoginRequest

// 	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
// 		h.logger.Error("failed decode request", zap.Error(err))
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 	}

// 	user, err := h.service.GetByLogin(r.Context(), input.Username, input.Email)
// 	if err != nil {
// 		h.logger.Error("failed get user", zap.Error(err))
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)

// 	if err = json.NewEncoder(w).Encode(dto.ToUserResponse(user)); err != nil {
// 		h.logger.Error("failed create response with user", zap.String("email", user.Email), zap.Error(err))
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
