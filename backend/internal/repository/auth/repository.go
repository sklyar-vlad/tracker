package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	model "github.com/sklyar-vlad/selfDev/internal/model/auth"
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

func (r *repository) CreateRefreshToken(ctx context.Context, tokens model.Tokens) error {
	query := `
	INSERT INTO refresh_tokens (token_hash, user_id, expires_at)
	VAlUES ($1, $2, $3)	
	`

	_, err := r.pool.Exec(ctx, query, tokens.RefreshToken, tokens.UserId, tokens.ExpiresAt)
	if err != nil {
		r.logger.Error("failed insert refresh token in database", zap.Error(err))
		return fmt.Errorf("failed insert refresh token in database: %v", err)
	}

	return nil
}

// func (r *repository) GetRefreshToken(ctx context.Context, userId uuid.UUID) (model.Tokens, error) {
// 	query := `
// 	SELECT token_hash, user_id, expires_at
// 	FROM refresh_tokens
// 	WHERE user_id = $1
// 	ORDER BY created_at DESC
// 	LIMIT 1
// 	`

// 	var refreshToken model.Tokens

// 	err := r.pool.QueryRow(ctx, query, userId).Scan(
// 		&refreshToken.TokenHash,
// 		&refreshToken.UserId,
// 		&refreshToken.ExpiresAt,
// 	)

// 	if errors.Is(err, pgx.ErrNoRows) {
// 		r.logger.Error("user not found", zap.Error(customErrors.ErrUserNotFound))
// 		return model.Tokens{}, customErrors.ErrUserNotFound
// 	}

// 	if err != nil {
// 		r.logger.Error("failed get refresh token by user id", zap.Error(err))
// 		return model.Tokens{}, err
// 	}

// 	r.logger.Info("success select token by user_id", zap.String("user_id", refreshToken.UserId.String()))

// 	return refreshToken, nil
// }
