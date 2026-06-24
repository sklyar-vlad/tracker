package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/sklyar-vlad/selfDev/database"
	"github.com/sklyar-vlad/selfDev/internal/config"
	"github.com/sklyar-vlad/selfDev/internal/handler"
	authHand "github.com/sklyar-vlad/selfDev/internal/handler/auth"
	userHand "github.com/sklyar-vlad/selfDev/internal/handler/user"
	authRepo "github.com/sklyar-vlad/selfDev/internal/repository/auth"
	userRepo "github.com/sklyar-vlad/selfDev/internal/repository/user"
	authSrv "github.com/sklyar-vlad/selfDev/internal/service/auth"
	userSrv "github.com/sklyar-vlad/selfDev/internal/service/user"
	"github.com/sklyar-vlad/selfDev/logger"
	"github.com/sklyar-vlad/selfDev/middleware"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed create config: %v", err)
	}

	logger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal("failed create logger", zap.Error(err))
	}

	defer func() {
		_ = logger.Sync()
	}()

	ctx := context.Background()
	pool, err := database.NewPostgres(ctx, cfg.Database)
	if err != nil {
		logger.Fatal("failed connect to the database", zap.Error(err))
	}

	defer pool.Close()

	authRepository := authRepo.NewRepository(pool, logger)
	userRepository := userRepo.NewRepository(pool, logger)

	userService := userSrv.NewService(userRepository, logger)
	authService := authSrv.NewService(authRepository, userService, cfg.JWT, logger)

	authHandler := authHand.NewHandler(authService, logger)
	userHandler := userHand.NewHandler(userService, logger)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, userHandler, authHandler)
	wrapped := middleware.CORS(mux)

	service := &http.Server{
		Addr:         ":8080",
		Handler:      wrapped,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("service started at port 8080.")

		if err := service.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := service.Shutdown(ctxShutdown); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
	} else {
		logger.Info("server stopped gracefully")
	}
}
