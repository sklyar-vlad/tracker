package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	appErrors "github.com/sklyar-vlad/selfDev/internal/errors"
	model "github.com/sklyar-vlad/selfDev/internal/model/auth"
)

type repository struct {
	pool   *pgxpool.Pool
	redis  *redis.Client
	logger *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, redis *redis.Client, logger *zap.Logger) *repository {
	return &repository{
		pool:   pool,
		redis:  redis,
		logger: logger,
	}
}

func (r *repository) CreateRefreshToken(ctx context.Context, tokens model.Tokens) error {
	query := `
	INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
	VAlUES ($1, $2, $3)	
	`

	r.logger.Info("get tokens", zap.String("refresh token", tokens.RefreshToken))
	_, err := r.pool.Exec(ctx, query, tokens.RefreshToken, tokens.UserId, tokens.ExpiresAt)
	if err != nil {
		r.logger.Error("failed insert refresh token in database", zap.Error(err))
		return fmt.Errorf("failed insert refresh token in database: %v", err)
	}

	return nil
}

func (r *repository) GetRefreshToken(ctx context.Context, userId uuid.UUID) (model.Tokens, error) {
	query := `
	SELECT token_hash, expires_at
	FROM refresh_tokens
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`

	var refreshToken model.Tokens

	err := r.pool.QueryRow(ctx, query, userId).Scan(
		&refreshToken.RefreshToken,
		&refreshToken.ExpiresAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("user not found", zap.Error(appErrors.ErrUserNotFound))
		return model.Tokens{}, appErrors.ErrUserNotFound
	}

	if err != nil {
		r.logger.Error("failed get refresh token by user id", zap.Error(err))
		return model.Tokens{}, err
	}

	r.logger.Info("success select token by user_id", zap.String("refresh token", refreshToken.RefreshToken))

	return refreshToken, nil
}

func (r *repository) DeleteRefreshToken(ctx context.Context, userId uuid.UUID) error {
	query := `
	DELETE FROM refresh_tokens
	WHERE user_id = $1`

	_, err := r.pool.Exec(ctx, query, userId)
	if err != nil {
		r.logger.Error("failed delete refresh token", zap.Error(err))
		return fmt.Errorf("failed delete refresh token: %v", err)
	}

	return nil
}

func (r *repository) SaveTokenVerify(ctx context.Context, token, userId string) error {
	key := "verify_email:" + token
	return r.redis.Set(ctx, key, userId, 5*time.Hour).Err()
}

func (r *repository) ConsumeToken(ctx context.Context, token string) (string, error) {
	key := "verify_email:" + token

	userID, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	_ = r.redis.Del(ctx, key)

	return userID, nil
}
