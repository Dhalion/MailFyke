package main

import (
	"net/http"
	"os"

	"github.com/Dhalion/MailFyke/internal/api"
	"github.com/Dhalion/MailFyke/internal/config"
	"github.com/Dhalion/MailFyke/internal/database"
	"github.com/Dhalion/MailFyke/internal/handler"
	"github.com/Dhalion/MailFyke/internal/httputil"
	"github.com/Dhalion/MailFyke/internal/middleware"
	"github.com/Dhalion/MailFyke/internal/smtp"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func runServer(cfg *config.Config) error {
	pool, err := database.NewPool(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	smtpServer := smtp.New(pool, cfg)
	err = smtpServer.Start(cfg.SMTPListenAddr)
	if err != nil {
		return err
	}
	defer func(smtpServer *smtp.Server) {
		err := smtpServer.Stop()
		if err != nil {

		}
	}(smtpServer)

	h := handler.New(pool, cfg)
	w := api.ServerInterfaceWrapper{
		Handler: h,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			httputil.WriteError(w, http.StatusBadRequest, err.Error())
		},
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/api/healthz", w.Health)
	r.Post("/api/auth/login", w.AuthLogin)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))

		r.Get("/api/users/me/organizations", w.ListOrganizations)

		r.Route("/api/organizations/{orgId}", func(r chi.Router) {
			r.Use(middleware.RequireOrgMembership)
			r.Get("/emails", w.ListEmails)
			r.Get("/emails/unread-count", w.UnreadCount)
			r.Get("/emails/{id}", w.GetEmail)
			r.Delete("/emails/{id}", w.DeleteEmail)
			r.Put("/emails/{id}/read", w.MarkRead)
		})
	})

	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	log.Info().Str("addr", cfg.ListenAddr).Msg("starting server")
	return http.ListenAndServe(cfg.ListenAddr, r)
}
