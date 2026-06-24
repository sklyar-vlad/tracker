package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/sklyar-vlad/selfDev/internal/config"
	appErrors "github.com/sklyar-vlad/selfDev/internal/errors"
	authModel "github.com/sklyar-vlad/selfDev/internal/model/auth"
	userModel "github.com/sklyar-vlad/selfDev/internal/model/user"
)

type UserService interface {
	CreateUser(ctx context.Context, username, email, password string) (userModel.User, error)
}

type Repository interface {
	CreateRefreshToken(ctx context.Context, Tokens authModel.Tokens) error
	// GetTokens(ctx context.Context, userId uuid.UUID) (authModel.Tokens, error)
	// DeleteTokens(ctx context.Context, userId uuid.UUID) error
}

type Service struct {
	repo        Repository
	userService UserService
	cfg         config.ConfigJWT
	logger      *zap.Logger
}

func NewService(repo Repository, userService UserService, config config.ConfigJWT, logger *zap.Logger) *Service {
	return &Service{repo: repo, userService: userService, cfg: config, logger: logger}
}

func (s *Service) Register(ctx context.Context, username, email, password string) (authModel.Tokens, error) {
	user, err := s.userService.CreateUser(ctx, username, email, password)

	if errors.Is(err, appErrors.ErrEmailAlreadyExists) {
		s.logger.Error("email already exists", zap.Error(err))
		return authModel.Tokens{}, appErrors.ErrEmailAlreadyExists
	}

	if errors.Is(err, appErrors.ErrUsernameAlreadyExists) {
		s.logger.Error("username is unvailable", zap.Error(err))
		return authModel.Tokens{}, appErrors.ErrUsernameAlreadyExists
	}

	if err != nil {
		s.logger.Error("failed create user", zap.Error(err))
		return authModel.Tokens{}, fmt.Errorf("failed create user: %v", err)
	}

	refreshTokenString, err := bcrypt.GenerateFromPassword([]byte(uuid.NewString()), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed hash generation", zap.Error(err))
		return authModel.Tokens{}, fmt.Errorf("failed hash generation: %v", err)
	}

	accessTokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, authModel.NewClaims(user.UserId)).
		SignedString([]byte(s.cfg.Secret))
	if err != nil {
		s.logger.Error("failed signed token", zap.Error(err))
		return authModel.Tokens{}, fmt.Errorf("failed signed token: %v", err)
	}

	var tokens authModel.Tokens
	tokens.AccessToken = accessTokenString
	tokens.RefreshToken = string(refreshTokenString)
	tokens.ExpiresAt = time.Now().AddDate(0, 1, 0)
	tokens.UserId = user.UserId

	err = s.repo.CreateRefreshToken(ctx, tokens)
	if err != nil {
		s.logger.Error("failed create refresh token", zap.Error(err))
		return authModel.Tokens{}, fmt.Errorf("failed create refresh token: %v", err)
	}

	s.logger.Info("success registered", zap.String("email", user.Email))
	return tokens, nil
}

// func (s *Service) Login(ctx context.Context, username, email, password string) (string, string, error) {
// 	user, err := s.repo.GetUserByEmail(ctx, email)

// 	if errors.Is(err, customErrors.ErrUserNotFound) {
// 		s.logger.Error("user not found", zap.Error(customErrors.ErrUserNotFound))
// 		return authModel.Tokens{}, authModel.Tokens{}, customErrors.ErrUserNotFound
// 	}

// 	if err != nil {
// 		s.logger.Error("failed select user")
// 		return authModel.Tokens{}, authModel.Tokens{}, err
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
// 	if err != nil {
// 		s.logger.Info("incorrect password", zap.Error(customErrors.ErrUnauthorized))
// 		return authModel.Tokens{}, authModel.Tokens{}, customErrors.ErrInvalidPassword
// 	}

// 	var Tokens model.Tokens
// 	TokensBytes := uuid.NewString()
// 	TokensHash, err := bcrypt.GenerateFromPassword([]byte(TokensBytes), bcrypt.DefaultCost)
// 	if err != nil {
// 		s.logger.Error("password hash generation failed", zap.String("email", email), zap.Error(err))
// 		return authModel.Tokens{}, authModel.Tokens{}, errors.New("invalid generate hash of password")
// 	}

// 	Tokens.TokenHash = string(TokensHash)
// 	Tokens.ExpiresAt = time.Now().AddDate(0, 1, 0)
// 	Tokens.UserId = user.UserId

// 	_, err = s.repo.CreateTokens(ctx, Tokens)
// 	if err != nil {
// 		s.logger.Error("failed create refresh token", zap.String("email", email), zap.Error(err))
// 		return authModel.Tokens{}, authModel.Tokens{}, errors.New("invalid create refresh token")
// 	}

// 	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, model.NewClaims(user.UserId)).
// 		SignedString([]byte("meow"))
// 	if err != nil {
// 		s.logger.Error("failed signed token", zap.String("email", email), zap.Error(err))
// 		return authModel.Tokens{}, authModel.Tokens{}, errors.New("invalid generate jwt")
// 	}

// 	s.logger.Info("success login", zap.String("email", user.Email))
// 	return accessToken, TokensBytes, nil
// }

// func (s *Service) Refresh(ctx context.Context, accessToken, Tokens string) (string, error) {
// 	token, err := jwt.ParseWithClaims(
// 		accessToken,
// 		&model.Claims{},
// 		func(t *jwt.Token) (interface{}, error) {
// 			return []byte("meow"), nil
// 		},
// 	)

// 	s.logger.Info("access token", zap.String("access token", accessToken))

// 	if err != nil {
// 		s.logger.Error("invalid validate jwt token", zap.Error(err))
// 		return authModel.Tokens{}, err
// 	}

// 	claims, ok := token.Claims.(*model.Claims)
// 	if !ok || !token.Valid {
// 		return authModel.Tokens{}, errors.New("invalid token")
// 	}

// 	TokensDB, err := s.repo.GetTokens(ctx, claims.UserId)
// 	if err != nil {
// 		return authModel.Tokens{}, err
// 	}

// 	if err = bcrypt.CompareHashAndPassword([]byte(TokensDB.TokenHash), []byte(Tokens)); err != nil {
// 		s.logger.Error("invalid refresh token", zap.Error(err))
// 		return authModel.Tokens{}, err
// 	}

// 	if time.Now().After(TokensDB.ExpiresAt) {
// 		s.logger.Error("refresh token expired")
// 		return authModel.Tokens{}, nil
// 	}

// 	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, model.NewClaims(claims.UserId)).
// 		SignedString([]byte("meow"))
// 	if err != nil {
// 		s.logger.Error("failed signed token", zap.String("user_id", claims.UserId.String()), zap.Error(err))
// 		return authModel.Tokens{}, errors.New("invalid generate jwt")
// 	}

// 	return newAccessToken, nil
// }
