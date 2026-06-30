package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	UpdateUser(ctx context.Context, user userModel.User) error
	GetByLogin(ctx context.Context, username, password string) (userModel.User, error)
	GetById(ctx context.Context, id uuid.UUID) (userModel.User, error)
}

type EmailAdapter interface {
	SendEmailVerification(email, token string) error
}

type Repository interface {
	CreateRefreshToken(ctx context.Context, Tokens authModel.Tokens) error
	GetRefreshToken(ctx context.Context, userId uuid.UUID) (authModel.Tokens, error)
	DeleteRefreshToken(ctx context.Context, userId uuid.UUID) error
	SaveTokenVerify(ctx context.Context, token, userId string) error
	ConsumeToken(ctx context.Context, token string) (string, error)
}

type Service struct {
	repo         Repository
	userService  UserService
	emailAdapter EmailAdapter
	cfg          config.ConfigJWT
	logger       *zap.Logger
}

func NewService(
	repo Repository,
	userService UserService,
	emailAdapter EmailAdapter,
	config config.ConfigJWT,
	logger *zap.Logger,
) *Service {
	return &Service{repo: repo, userService: userService, emailAdapter: emailAdapter, cfg: config, logger: logger}
}

func (s *Service) Register(ctx context.Context, username, email, password string) error {
	user, err := s.userService.CreateUser(ctx, username, email, password)

	if errors.Is(err, appErrors.ErrEmailAlreadyExists) {
		s.logger.Error("email already exists", zap.Error(err))
		return appErrors.ErrEmailAlreadyExists
	}

	if errors.Is(err, appErrors.ErrUsernameAlreadyExists) {
		s.logger.Error("username is unvailable", zap.Error(err))
		return appErrors.ErrUsernameAlreadyExists
	}

	if err != nil {
		s.logger.Error("failed create user", zap.Error(err))
		return fmt.Errorf("failed create user: %v", err)
	}

	token, err := authModel.NewTokenVerify()
	if err != nil {
		s.logger.Error("failed create verify token", zap.Error(err))
		return fmt.Errorf("failed create verify token: %v", err)
	}

	err = s.repo.SaveTokenVerify(ctx, token.TokenVer, user.UserId.String())
	if err != nil {
		s.logger.Error("failed save verify token in redis", zap.Error(err))
		return fmt.Errorf("failed create save verify token in redis: %v", err)
	}

	go func() {
		err := s.emailAdapter.SendEmailVerification(email, token.TokenVer)
		if err != nil {
			s.logger.Error("failed send message for verification", zap.Error(err))
		}
	}()

	s.logger.Info("success registered", zap.String("email", user.Email))
	return nil
}

func (s *Service) Login(ctx context.Context, username, email, password string) (string, string, error) {
	user, err := s.userService.GetByLogin(ctx, username, email)

	if errors.Is(err, appErrors.ErrUserNotFound) {
		s.logger.Error("user not found", zap.Error(appErrors.ErrUserNotFound))
		return "", "", appErrors.ErrUserNotFound
	}

	if err != nil {
		s.logger.Error("failed get user", zap.Error(err))
		return "", "", err
	}

	if !user.EmailVerified {
		s.logger.Error("user's email not verified", zap.Error(appErrors.ErrEmailNotVerified))
		return "", "", appErrors.ErrEmailNotVerified
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", appErrors.ErrInvalidPassword
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, authModel.NewRefreshClaims(user.UserId)).SignedString([]byte(s.cfg.Secret))
	if err != nil {
		s.logger.Error("failed hash generation", zap.Error(err))
		return "", "", fmt.Errorf("failed hash generation: %v", err)
	}

	refreshTokenHash := sha256.Sum256([]byte(refreshToken))

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, authModel.NewAccessClaims(user.UserId)).SignedString([]byte(s.cfg.Secret))
	if err != nil {
		s.logger.Error("failed signed token", zap.Error(err))
		return "", "", fmt.Errorf("failed signed token: %v", err)
	}

	var tokens authModel.Tokens
	tokens.AccessToken = accessToken
	tokens.RefreshToken = hex.EncodeToString(refreshTokenHash[:])
	tokens.ExpiresAt = time.Now().AddDate(0, 1, 0)
	tokens.UserId = user.UserId

	s.logger.Info("create model of tokens", zap.String("refresh token", tokens.RefreshToken), zap.String("access token", tokens.AccessToken))

	err = s.repo.CreateRefreshToken(ctx, tokens)
	if err != nil {
		s.logger.Error("failed create refresh token", zap.Error(err))
		return "", "", fmt.Errorf("failed create refresh token: %v", err)
	}

	s.logger.Info("success login", zap.String("email", user.Email))
	return refreshToken, accessToken, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	token, err := jwt.ParseWithClaims(
		refreshToken,
		&authModel.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.Secret), nil
		},
	)
	if err != nil {
		s.logger.Error("invalid validate jwt token", zap.Error(err))
		return err
	}

	claims, ok := token.Claims.(*authModel.Claims)
	if !ok || !token.Valid {
		return errors.New("invalid token")
	}

	err = s.repo.DeleteRefreshToken(ctx, claims.UserId)
	if err != nil {
		s.logger.Error("failed delete refresh token", zap.Error(err))
		return err
	}

	s.logger.Info("success logout user", zap.String("user_id", claims.UserId.String()))
	return nil
}



func (s *Service) ConfirmEmail(ctx context.Context, token string) error {
	userId, err := s.repo.ConsumeToken(ctx, token)
	if err != nil {
		return err
	}

	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	userEmailVerified, err := s.userService.GetById(ctx, userIdUUID)
	if err != nil {
		s.logger.Error("failed get user by id", zap.Error(err))
		return fmt.Errorf("failed get user by id: %v", err)
	}

	userEmailVerified.EmailVerified = true

	if err = s.userService.UpdateUser(ctx, userEmailVerified); err != nil {
		return fmt.Errorf("failed verified email: %v", err)
	}

	return nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, error) {
	token, err := jwt.ParseWithClaims(
		refreshToken,
		&authModel.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.cfg.Secret), nil
		},
	)
	if err != nil {
		s.logger.Error("invalid validate jwt token", zap.Error(err))
		return "", err
	}

	claims, ok := token.Claims.(*authModel.Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	refreshTokenInDB, err := s.repo.GetRefreshToken(ctx, claims.UserId)
	if err != nil {
		return "", err
	}

	refreshTokenHash := sha256.Sum256([]byte(refreshToken))

	if hex.EncodeToString(refreshTokenHash[:]) != refreshTokenInDB.RefreshToken {
		s.logger.Error("invalid refresh token", zap.Error(err))
		return "", err
	}

	if time.Now().After(refreshTokenInDB.ExpiresAt) {
		s.logger.Error("refresh token expired")
		return "", nil
	}

	newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, authModel.NewAccessClaims(claims.UserId)).SignedString([]byte(s.cfg.Secret))
	if err != nil {
		s.logger.Error("failed signed token", zap.String("user_id", claims.UserId.String()), zap.Error(err))
		return "", fmt.Errorf("failed signed token: %v", err)
	}

	return newAccessToken, nil
}
