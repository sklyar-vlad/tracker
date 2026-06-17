package user

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"github.com/google/uuid"

	customErrors "github.com/sklyar-vlad/selfDev/internal/errors"
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

func (r *repository) Register(ctx context.Context, user model.User) (model.User, error) {
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

func (r *repository) GetRefreshToken(ctx context.Context, user model.User) (model.RefreshToken, error) {
	query := `
	SELECT user_id, token_hash, created_at, expires_at
	FROM refresh_tokens
	WHERE user_id = $1 and expires_at > NOW();
	`

	var refreshToken model.RefreshToken

	err := r.pool.QueryRow(ctx, query, user.User_id).Scan(
		&refreshToken.User_id,
		&refreshToken.TokenHash,
		&refreshToken.CreatedAt,
		&refreshToken.ExpiresAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("refresh token was expired", zap.Error(customErrors.ErrTokenWasExpired))
		return model.RefreshToken{}, customErrors.ErrTokenWasExpired
	}

	if err != nil {
		r.logger.Error("failed get user by user_id", zap.Error(err))
		return model.RefreshToken{}, err
	}


	r.logger.Info("success select refresh token", zap.String("email", user.Email))

	return refreshToken, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `
	SELECT user_id, role, username, email, password
	FROM users
	WHERE email = $1
	`

	var user model.User

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.User_id,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("user not found", zap.Error(customErrors.ErrUserNotFound))
		return model.User{}, customErrors.ErrUserNotFound
	}

	if err != nil {
		r.logger.Error("failed get user by email", zap.Error(err))
		return model.User{}, err
	}

	r.logger.Info("success select user by email", zap.String("email", user.Email))
	return user, nil
}

func (r *repository) CreateToken(ctx context.Context, user_id uuid.UUID) (model.User, error) {
	query := `
	SELECT user_id, role, username, email, password
	FROM users
	WHERE email = $1
	`

	var user model.User

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.User_id,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("user not found", zap.Error(customErrors.ErrUserNotFound))
		return model.User{}, customErrors.ErrUserNotFound
	}

	if err != nil {
		r.logger.Error("failed get user by email", zap.Error(err))
		return model.User{}, err
	}

	r.logger.Info("success select user by email", zap.String("email", user.Email))
	return user, nil
}