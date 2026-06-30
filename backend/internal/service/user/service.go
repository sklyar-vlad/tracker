package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	appErrors "github.com/sklyar-vlad/selfDev/internal/errors"
	model "github.com/sklyar-vlad/selfDev/internal/model/user"
)

type Repository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	Update(ctx context.Context, user model.User) error
	GetByLogin(ctx context.Context, login string) (model.User, error)
	GetById(ctx context.Context, userId uuid.UUID) (model.User, error)
	// Update(ctx context.Context, user model.User) (model.User, error)
	// Delete(ctx context.Context, user model.User) error
}

type Service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) CreateUser(ctx context.Context, username, email, password string) (model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed password hash generation", zap.String("email", email), zap.Error(err))
		return model.User{}, fmt.Errorf("failed password hash generation: %v", err)
	}

	user, err := model.NewUser(username, email, string(passwordHash))

	if errors.Is(err, appErrors.ErrInvalidEmail) {
		s.logger.Error("invalid email", zap.Error(err))
		return model.User{}, appErrors.ErrInvalidEmail
	}

	if errors.Is(err, appErrors.ErrInvalidPassword) {
		s.logger.Error("invalid password", zap.Error(err))
		return model.User{}, appErrors.ErrInvalidPassword
	}

	if err != nil {
		s.logger.Error("failed create user model", zap.String("email", email), zap.Error(err))
		return model.User{}, fmt.Errorf("failed create user model: %v", err)
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, user model.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetByLogin(ctx context.Context, username, email string) (model.User, error) {
	var login string
	if email == "" {
		login = username
	} else {
		login = email
	}

	user, err := s.repo.GetByLogin(ctx, login)

	if errors.Is(err, appErrors.ErrUserNotFound) {
		s.logger.Error("user not found", zap.Error(err))
		return model.User{}, appErrors.ErrUserNotFound
	}

	if err != nil {
		s.logger.Error("failed get user", zap.String("email/username", login), zap.Error(err))
		return model.User{}, fmt.Errorf("failed get user: %v", err)
	}

	return user, nil
}

func (s *Service) GetById(ctx context.Context, userId uuid.UUID) (model.User, error) {
	user, err := s.repo.GetById(ctx, userId)

	if errors.Is(err, appErrors.ErrUserNotFound) {
		s.logger.Error("user not found", zap.Error(err))
		return model.User{}, appErrors.ErrUserNotFound
	}

	if err != nil {
		s.logger.Error("failed get user", zap.String("id", userId.String()), zap.Error(err))
		return model.User{}, fmt.Errorf("failed get user: %v", err)
	}

	return user, nil
}
