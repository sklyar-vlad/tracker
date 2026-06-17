package user

import (
	"context"
	"errors"

	"github.com/sklyar-vlad/selfDev/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	// GetByEmail(email string) (model.User, error)
	// Update(user model.User) model.User
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, username, email, password string) (model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, errors.New("invalid generate hash of password")
	}

	user, err := model.NewUser(username, email, string(passwordHash))
	if err != nil {
		return model.User{}, errors.New("invalid create user")
	}

	return s.repo.Create(ctx, user)
}
