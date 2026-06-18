package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	DatabaseURL       string
	ListenAddr        string
	SMTPListenAddr    string
	SMTPMaxMailSize   int
	SMTPMaxConns      int
	SMTPAllowInsecure bool
	SMTPDomain        string
	SMTPDebug         bool
	LogLevel          string
	JWTSecret         string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://mailfyke:mailfyke@localhost:5432/mailfyke?sslmode=disable"),
		ListenAddr:        getEnv("LISTEN_ADDR", ":5789"),
		SMTPListenAddr:    getEnv("SMTP_LISTEN_ADDR", ":2525"),
		SMTPMaxMailSize:   getEnvInt("SMTP_MAX_MAIL_SIZE", 26214400),
		SMTPMaxConns:      getEnvInt("SMTP_MAX_CONNECTIONS", 100),
		SMTPAllowInsecure: getEnvBool("SMTP_ALLOW_INSECURE", false),
		SMTPDomain:        getEnv("SMTP_DOMAIN", "localhost"),
		SMTPDebug:         getEnvBool("SMTP_DEBUG", false),
		LogLevel:          getEnv("LOG_LEVEL", "debug"),
		JWTSecret:         getEnv("JWT_SECRET", "dev-secret-change-in-production"),
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

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}

	return b
}

func initLogger(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(lvl)
}
