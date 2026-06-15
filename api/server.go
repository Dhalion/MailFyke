package main

import (
	"github.com/chris/MailFyke/internal/config"
	"github.com/chris/MailFyke/internal/database"
	"github.com/chris/MailFyke/internal/handler"
	"github.com/chris/MailFyke/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
)

func runServer(cfg *config.Config) error {
	pool, err := database.NewPool(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(middleware.Auth(cfg.JWTSecret))

	h := handler.New(pool, cfg)
	r.Get("/api/healthz", h.Health)

	// TODO: mount SMTP server

	zerolog.Info().Str("addr", cfg.ListenAddr).Msg("starting server")
	return http.ListenAndServe(cfg.ListenAddr, r)
}
