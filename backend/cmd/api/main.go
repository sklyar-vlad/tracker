package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	"github.com/sklyar-vlad/selfDev/database"
	"github.com/sklyar-vlad/selfDev/internal/handler"
	"github.com/sklyar-vlad/selfDev/internal/handler/user"
	userRepo "github.com/sklyar-vlad/selfDev/internal/repository/user"
	userSrv "github.com/sklyar-vlad/selfDev/internal/service/user"
	"github.com/sklyar-vlad/selfDev/logger"
	"github.com/sklyar-vlad/selfDev/middleware"
	_ "github.com/sklyar-vlad/selfDev/swagger"
)

//	@title			Swagger SelfDev API
//	@version		1.0
//	@description	This is a server of self-dev tracker.

//	@contact.name	API Support
//	@contact.url	t.me/sklyarvlad
//	@contact.email	sklyarvladislavtl@gmail.com

// @host		localhost:8080
// @BasePath	/api
func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = logger.Sync()
	}()

	if err := godotenv.Load(); err != nil {
		logger.Fatal("invalid load .env file:", zap.Error(err))
	}

	ctx := context.Background()
	pool, err := database.NewPostgres(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("invalid connect with database", zap.Error(err))
	}

	defer pool.Close()

	userRepository := userRepo.NewRepository(pool, logger)
	userService := userSrv.NewService(userRepository, logger)
	userHandler := user.NewHandler(userService, logger)

	mux := http.NewServeMux()
	mux.Handle("GET /api/swagger/", httpSwagger.WrapHandler)
	handler.RegisterRoutes(mux, userHandler)
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
