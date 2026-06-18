package user

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	customErrors "github.com/sklyar-vlad/selfDev/internal/errors"
	"github.com/sklyar-vlad/selfDev/internal/model"
)

type Repository interface {
	Register(ctx context.Context, user model.User) (model.User, error)
	Auth(ctx context.Context, user model.User) error
	GetRefreshToken(ctx context.Context, user model.User) (model.RefreshToken, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
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
		s.logger.Info("incorrect password or email", zap.Error(customErrors.ErrUnauthorized))
		return "", "", customErrors.ErrInvalidPassword
	}

	s.logger.Info("success authorizated", zap.String("email", user.Email))
	return
}

func (s *Service) Auth(ctx context.Context, username, email, password string) (string, string, error) {

	user, err := s.repo.GetUserByEmail(ctx, email)

	if errors.Is(err, customErrors.ErrUserNotFound) {
		s.logger.Error("user not found", zap.Error(custors.ErrUserNotFound))
		return "", "", customErrors.ErrUserNotFound
	}

	if err != nil {
		s.logger.Error("failed select user")
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		s.logger.Info("incorrect password or email", zap.Error(customErrors.ErrUnauthorized))
		return "", "", customErrors.ErrInvalidPassword
	}

	refreshToken, err := s.repo.GetRefreshToken(ctx, user)

	if errors.Is(err, customErrors.ErrTokenWasExpired) {
		refreshToken = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		s.repo.CreateRefreshToken(ctx, user.User_id)
	}

	if err != nil {
		s.logger.Error("failed select refresh token")
		return "", "", err
	}

	// accessToken :=

	s.logger.Info("success authorizated", zap.String("email", user.Email))
	return access_token, refreshToken.TokenHash, nil
}
