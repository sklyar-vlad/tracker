package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigLogger struct {
	Env string
}

type ConfigServer struct {
	Host         string
	Addr         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type ConfigDatabase struct {
	PostgresURL string
	RedisURL    string
}

type ConfigJWT struct {
	Secret string
}

type ConfigEmailSender struct {
	ApiKey   string
	From     string
	Subject  string
	Endpoint string
}

type config struct {
	Server      ConfigServer
	Logger      ConfigLogger
	Database    ConfigDatabase
	JWT         ConfigJWT
	EmailSender ConfigEmailSender
}

func NewConfig() (config, error) {
	if err := godotenv.Load(); err != nil {
		return config{}, err
	}

	readTimeout, _ := strconv.Atoi(getEnv("READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnv("WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnv("IDLE_TIMEOUT", "60"))

	cfg := config{
		Server: ConfigServer{
			Host:         getEnv("HOST", "localhost"),
			Addr:         getEnv("ADDR", ":8080"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		Logger: ConfigLogger{
			Env: getEnv("ENV", "production"),
		},
		Database: ConfigDatabase{
			PostgresURL: getEnv("POSTGRES_URL", "postgres://admin:admin@db:5432/self-dev"),
			RedisURL:    getEnv("REDIS_URL", "redis://:admin@redis:6379/0"),
		},
		JWT: ConfigJWT{
			Secret: getEnv("SECRET", ""),
		},
		EmailSender: ConfigEmailSender{
			ApiKey:   getEnv("RESEND_API_KEY", ""),
			From:     getEnv("FROM_EMAIL", ""),
			Subject:  getEnv("SUBJECT", ""),
			Endpoint: getEnv("ENDPOINT", ""),
		},
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}
