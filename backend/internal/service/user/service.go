package user

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/sklyar-vlad/selfDev/internal/model"
)

type Repository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
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
	return s.repo.Create(ctx, user)
}
