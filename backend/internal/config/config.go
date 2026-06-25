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
	Addr         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type ConfigDatabase struct {
	DatabaseURL string
}

type ConfigJWT struct {
	Secret string
}

type ConfigEmailSender struct {
	ApiKey  string
	From    string
	Html    string
	Subject string
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
			Addr:         getEnv("ADDR", ":8080"),
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		Logger: ConfigLogger{
			Env: getEnv("ENV", "production"),
		},
		Database: ConfigDatabase{
			DatabaseURL: getEnv("DATABASE_URL", "postgres://admin:admin@db:5432/self-dev"),
		},
		JWT: ConfigJWT{
			Secret: getEnv("SECRET", ""),
		},
		EmailSender: ConfigEmailSender{
			ApiKey:  getEnv("RESEND_API_KEY", ""),
			From:    getEnv("FROM_EMAIL", ""),
			Html:    getEnv("HTML", ""),
			Subject: getEnv("SUBJECT", ""),
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
