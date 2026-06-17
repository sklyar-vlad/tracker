package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/sklyar-vlad/selfDev/internal/model"
)

type repository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, logger *zap.Logger) *repository {
	return &repository{
		pool:   pool,
		logger: logger,
	}
}

func (r *repository) Create(ctx context.Context, user model.User) (model.User, error) {
	query := `
	INSERT INTO users (role, username, email, password)
	VALUES ($1, $2, $3,	$4)
	`

	_, err := r.pool.Exec(ctx, query, user.Role, user.Username, user.Email, user.Password)
	if err != nil {
		r.logger.Error("invalid insert user into database", zap.Error(err))
		return model.User{}, err
	}

	r.logger.Info("success insert user in database", zap.String("email", user.Email))
	return user, nil
}
