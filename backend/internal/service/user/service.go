package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	customErrors "github.com/sklyar-vlad/selfDev/internal/errors"
	"github.com/sklyar-vlad/selfDev/internal/model"
)

type Repository interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	CreateRefreshToken(ctx context.Context, refreshToken model.RefreshToken) (model.RefreshToken, error)
	GetRefreshToken(ctx context.Context, userId uuid.UUID) (model.RefreshToken, error)
}

type Service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Register(ctx context.Context, username, email, password string) (model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("password hash generation failed", zap.String("email", email), zap.Error(err))
		return model.User{}, errors.New("invalid generate hash of password")
	}

	user, err := model.NewUser(username, email, string(passwordHash))
	if err != nil {
		s.logger.Error("create model user failed", zap.String("email", email), zap.Error(err))
		return model.User{}, errors.New("invalid create user")
	}

	s.logger.Info("success created model user", zap.String("email", user.Email))
	return s.repo.Register(ctx, user)
}

func (s *Service) Login(ctx context.Context, username, email, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)

	if errors.Is(err, customErrors.ErrUserNotFound) {
		s.logger.Error("user not found", zap.Error(customErrors.ErrUserNotFound))
		return "", "", customErrors.ErrUserNotFound
	}

	if err != nil {
		s.logger.Error("failed select user")
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.logger.Info("incorrect password", zap.Error(customErrors.ErrUnauthorized))
		return "", "", customErrors.ErrInvalidPassword
	}

	var refreshToken model.RefreshToken
	refreshTokenBytes := uuid.NewString()
	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshTokenBytes), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("password hash generation failed", zap.String("email", email), zap.Error(err))
		return "", "", errors.New("invalid generate hash of password")
	}

	refreshToken.TokenHash = string(refreshTokenHash)
	refreshToken.ExpiresAt = time.Now().AddDate(0, 1, 0)
	refreshToken.UserId = user.UserId

	_, err = s.repo.CreateRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Error("failed create refresh token", zap.String("email", email), zap.Error(err))
		return "", "", errors.New("invalid create refresh token")
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, model.NewClaims(user.UserId)).
		SignedString([]byte("meow"))
	if err != nil {
		s.logger.Error("failed signed token", zap.String("email", email), zap.Error(err))
		return "", "", errors.New("invalid generate jwt")
	}

	s.logger.Info("success login", zap.String("email", user.Email))
	return accessToken, refreshTokenBytes, nil
}

func (s *Service) Refresh(ctx context.Context, accessToken, refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&model.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte("meow"), nil
		},
	)

	s.logger.Info("access token", zap.String("access token", accessToken))

	if err != nil {
		s.logger.Error("invalid validate jwt token", zap.Error(err))
		return "", err
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	refreshTokenDB, err := s.repo.GetRefreshToken(ctx, claims.UserId)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(refreshTokenDB.TokenHash), []byte(refreshToken)); err != nil {
		s.logger.Error("invalid refresh token", zap.Error(err))
		return "", err
	}

	if time.Now().After(refreshTokenDB.ExpiresAt) {
		s.logger.Error("refresh token expired")
		return "", nil
	}

	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, model.NewClaims(claims.UserId)).
		SignedString([]byte("meow"))
	if err != nil {
		s.logger.Error("failed signed token", zap.String("user_id", claims.UserId.String()), zap.Error(err))
		return "", errors.New("invalid generate jwt")
	}

	return newAccessToken, nil
}
