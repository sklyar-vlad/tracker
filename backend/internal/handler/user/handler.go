package user

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/sklyar-vlad/selfDev/internal/errors"
	"github.com/sklyar-vlad/selfDev/internal/handler/user/dto"
	"github.com/sklyar-vlad/selfDev/internal/model"
)

type Service interface {
	Register(ctx context.Context, username, email, password string) (model.User, error)
	Login(ctx context.Context, username, email, password string) (string, string, error)
	Refresh(ctx context.Context, accessToken, refreshToken string) (string, error)
}

type handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *handler {
	return &handler{service: service, logger: logger}
}

// Register godoc
//
//	@Summary		Register new user
//	@Description	Creates a new user account
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateUserRequest	true	"User data"
//	@Success		201		{object}	dto.UserResponse
//	@Failure		400		{string}	string
//	@Failure		500		{string}	string
//	@Router			/register [post]
func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("decode request failed", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if input.Email == "" {
		h.logger.Error("invalid email address", zap.Error(errors.ErrInvalidEmail))
		http.Error(w, errors.ErrInvalidEmail.Error(), http.StatusInternalServerError)
		return
	}

	if len(input.Password) < 6 {
		h.logger.Error("invalid password, it should be more than 6 symbols", zap.Error(errors.ErrInvalidPassword))
		http.Error(w, errors.ErrInvalidPassword.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.service.Register(r.Context(), input.Username, input.Email, input.Password)
	if err != nil {
		h.logger.Error("failed create user", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(dto.ToUserResponse(user)); err != nil {
		h.logger.Error("failed create response with user model", zap.String("email", user.Email), zap.Error(err))
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("decode request failed", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if input.Email == "" {
		h.logger.Error("invalid email address", zap.Error(errors.ErrInvalidEmail))
		http.Error(w, errors.ErrInvalidEmail.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := h.service.Login(r.Context(), input.Username, input.Email, input.Password)
	if err != nil {
		h.logger.Error("failed authorization", zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   12 * 60 * 60,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   30 * 24 * 60 * 60,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input dto.TokenRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Error("decode request failed", zap.Error(err))
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	h.logger.Info("access token", zap.String("access token", input.AccessToken))
	accessToken, err := h.service.Refresh(r.Context(), input.AccessToken, input.RefreshToken)
	if err != nil {
		h.logger.Error("failed authorization", zap.Error(err))
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   12 * 60 * 60,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
