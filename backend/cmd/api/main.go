package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sklyar-vlad/selfDev/database"
	"github.com/sklyar-vlad/selfDev/internal/handler"
	"github.com/sklyar-vlad/selfDev/internal/handler/user"
	userRepo "github.com/sklyar-vlad/selfDev/internal/repository/user"
	userSrv "github.com/sklyar-vlad/selfDev/internal/service/user"
	"github.com/sklyar-vlad/selfDev/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := database.NewPostgres(
		ctx,
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepository := userRepo.NewRepository(pool)
	userService := userSrv.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux, userHandler)
	wrapped := middleware.CORS(mux)

	service := &http.Server{
		Addr:         ":8080",
		Handler:      wrapped,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("service started at localhost:8080")
	log.Fatal(service.ListenAndServe())
}
