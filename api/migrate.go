package main

import (
	"context"
	"os"

	"github.com/Dhalion/MailFyke/internal/config"
	"github.com/Dhalion/MailFyke/internal/database"
	"github.com/Dhalion/MailFyke/internal/migrations"
	"github.com/rs/zerolog"
)

func runMigrate(cfg *config.Config) error {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	ctx := context.Background()

	pool, err := database.NewPool(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	up, err := migrations.FS.ReadFile("000001_initial.up.sql")
	if err != nil {
		return err
	}

	if _, err := pool.Exec(ctx, string(up)); err != nil {
		return err
	}

	log.Info().Msg("migrations complete")
	return nil
}
