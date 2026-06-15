package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL       string
	ListenAddr        string
	SMTPListenAddr    string
	SMTPMaxMailSize   int
	SMTPMaxConns      int
	LogLevel          string
	JWTSecret         string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://mailfyke:mailfyke@localhost:5432/mailfyke?sslmode=disable"),
		ListenAddr:     getEnv("LISTEN_ADDR", ":5789"),
		SMTPListenAddr: getEnv("SMTP_LISTEN_ADDR", ":2525"),
		SMTPMaxMailSize: getEnvInt("SMTP_MAX_MAIL_SIZE", 26214400),
		SMTPMaxConns:   getEnvInt("SMTP_MAX_CONNECTIONS", 100),
		LogLevel:       getEnv("LOG_LEVEL", "debug"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-in-production"),
	}

	initLogger(cfg.LogLevel)
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return fallback
}

func initLogger(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)
}
