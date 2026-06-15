package main

import (
	"github.com/chris/MailFyke/internal/config"
	"github.com/rs/zerolog"
	"os"
)

func runMigrate(cfg *config.Config) error {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log.Info().Str("url", cfg.DatabaseURL).Msg("running migrations")
	// TODO: implement golang-migrate
	return nil
}
