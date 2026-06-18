package smtp

import (
	"github.com/Dhalion/MailFyke/internal/config"
	"github.com/Dhalion/MailFyke/internal/database/queries"
	"github.com/emersion/go-smtp"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Backend struct {
	pool   *pgxpool.Pool
	q      *queries.Queries
	config *config.Config
}

func (be *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{
		backend: be,
		auth:    false,
		config:  be.config,
		pool:    be.pool,
		q:       be.q,
	}, nil
}
