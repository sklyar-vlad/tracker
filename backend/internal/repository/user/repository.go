package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	appErrors "github.com/sklyar-vlad/selfDev/internal/errors"
	model "github.com/sklyar-vlad/selfDev/internal/model/user"
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
	INSERT INTO users (user_id, role, username, email, email_verified, password)
	VALUES ($1, $2, $3,	$4, $5, $6)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		user.UserId,
		user.Role,
		user.Username,
		user.Email,
		user.EmailVerified,
		user.Password,
	)
	if err != nil {
		r.logger.Error("failed insert user in database", zap.Error(err))
		return model.User{}, mapDBError(err)
	}

	r.logger.Info("success insert user in database", zap.String("email", user.Email))

	return user, nil
}

func (r *repository) Update(ctx context.Context, user model.User) error {
	query := `
	UPDATE users
	SET 
		role = $2,
		username = $3,
		email = $4,
		email_verified = $5,
		password = $6
	WHERE user_id = $1
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		user.UserId,
		user.Role,
		user.Username,
		user.Email,
		user.EmailVerified,
		user.Password,
	)

	if err != nil {
		r.logger.Error("failed update user in database", zap.Error(err))
		return fmt.Errorf("failed update user in database: %v", err)
	}

	r.logger.Info("success update user in database", zap.String("email", user.Email))

	return nil
}

func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	query := `
	SELECT user_id, role, username, email, email_verified, password
	FROM users
	WHERE email = $1 or username = $1
	LIMIT 1
	`

	var user model.User

	err := r.pool.QueryRow(ctx, query, login).Scan(
		&user.UserId,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.EmailVerified,
		&user.Password,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("user not found", zap.Error(err))
		return model.User{}, appErrors.ErrUserNotFound
	}

	if err != nil {
		r.logger.Error("failed get user from database", zap.Error(err))
		return model.User{}, fmt.Errorf("failed get user from database: %v", err)
	}

	return user, nil
}

func (r *repository) GetById(ctx context.Context, userId uuid.UUID) (model.User, error) {
	query := `
	SELECT user_id, role, username, email, email_verified, password
	FROM users
	WHERE user_id = $1
	LIMIT 1
	`
	var user model.User

	err := r.pool.QueryRow(ctx, query, userId).Scan(
		&user.UserId,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.EmailVerified,
		&user.Password,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("user not found", zap.Error(err))
		return model.User{}, appErrors.ErrUserNotFound
	}

	if err != nil {
		r.logger.Error("failed get user from database", zap.Error(err))
		return model.User{}, fmt.Errorf("failed get user from database: %v", err)
	}

	return user, nil
}

func mapDBError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			switch pgErr.ConstraintName {
			case "users_email_key":
				return appErrors.ErrEmailAlreadyExists
			case "users_username_key":
				return appErrors.ErrUsernameAlreadyExists
			}
		}
	}

	return fmt.Errorf("database error: %w", err)
}

// func (r *repository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
// 	query := `
// 	SELECT user_id, role, username, email, password
// 	FROM users
// 	WHERE email = $1
// 	`

// 	var user model.User

// 	err := r.pool.QueryRow(ctx, query, email).Scan(
// 		&user.UserId,
// 		&user.Role,
// 		&user.Username,
// 		&user.Email,
// 		&user.Password,
// 	)

// 	if errors.Is(err, pgx.ErrNoRows) {
// 		r.logger.Error("user not found", zap.Error(customErrors.ErrUserNotFound))
// 		return model.User{}, customErrors.ErrUserNotFound
// 	}

// 	if err != nil {
// 		r.logger.Error("failed get user by email", zap.Error(err))
// 		return model.User{}, err
// 	}

// 	r.logger.Info("success select user by email", zap.String("email", user.Email))
// 	return user, nil
// }
