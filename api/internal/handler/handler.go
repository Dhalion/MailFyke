package handler

import (
	"github.com/chris/MailFyke/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
	cfg  *config.Config
}

func New(pool *pgxpool.Pool, cfg *config.Config) *Handler {
	return &Handler{
		pool: pool,
		cfg:  cfg,
	}
}
