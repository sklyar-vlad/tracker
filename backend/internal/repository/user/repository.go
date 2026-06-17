package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sklyar-vlad/selfDev/internal/model"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}

func (r *repository) Create(ctx context.Context, user model.User) (model.User, error) {
	query := `
	INSERT INTO users (role, username, email, password)
	VALUES ($1, $2, $3,	$4)
	`

	_, err := r.pool.Exec(ctx, query, user.Role, user.Username, user.Email, user.Password)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
