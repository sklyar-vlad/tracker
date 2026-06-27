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
	UpdateUser(ctx context.Context, user userModel.User) error
	GetByLogin(ctx context.Context, username, password string) (userModel.User, error)
	GetById(ctx context.Context, id uuid.UUID) (userModel.User, error)
}

type EmailAdapter interface {
	SendEmailVerification(email, token string) error
}

type Repository interface {
	CreateRefreshToken(ctx context.Context, Tokens authModel.Tokens) error
	SaveTokenVerify(ctx context.Context, token, userId string) error
	ConsumeToken(ctx context.Context, token string) (string, error)
	// GetTokens(ctx context.Context, userId uuid.UUID) (authModel.Tokens, error)
	// DeleteTokens(ctx context.Context, userId uuid.UUID) error
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

func (s *Service) Login(ctx context.Context, username, email, password string) (authModel.Tokens, error) {
	user, err := s.userService.GetByLogin(ctx, username, email)

	if errors.Is(err, appErrors.ErrUserNotFound) {
		s.logger.Error("user not found", zap.Error(appErrors.ErrUserNotFound))
		return authModel.Tokens{}, appErrors.ErrUserNotFound
	}

	if err != nil {
		s.logger.Error("failed get user", zap.Error(err))
		return authModel.Tokens{}, err
	}

	if !user.EmailVerified {
		s.logger.Error("user's email not verified", zap.Error(appErrors.ErrEmailNotVerified))
		return authModel.Tokens{}, appErrors.ErrEmailNotVerified
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return authModel.Tokens{}, appErrors.ErrInvalidPassword
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

	s.logger.Info("success login", zap.String("email", user.Email))
	return tokens, nil
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
	userEmailVerified.EmailVerified = true
	
	if err = s.userService.UpdateUser(ctx, userEmailVerified); err != nil {
		return fmt.Errorf("failed verified email: %v", err)
	}

	return nil	
}

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
