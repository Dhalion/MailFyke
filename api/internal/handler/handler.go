package handler

import (
	"github.com/Dhalion/MailFyke/internal/api"
	"github.com/Dhalion/MailFyke/internal/config"
	"github.com/Dhalion/MailFyke/internal/database/queries"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	api.Unimplemented
	pool    *pgxpool.Pool
	queries *queries.Queries
	cfg     *config.Config
}

func New(pool *pgxpool.Pool, cfg *config.Config) *Handler {
	return &Handler{
		pool:    pool,
		queries: queries.New(pool),
		cfg:     cfg,
	}
}
