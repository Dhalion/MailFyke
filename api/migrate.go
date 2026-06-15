package main

import (
	"github.com/chris/MailFyke/internal/config"
	"github.com/rs/zerolog"
)

func runMigrate(cfg *config.Config) error {
	zerolog.Info().Str("url", cfg.DatabaseURL).Msg("running migrations")
	// TODO: implement golang-migrate
	return nil
}
