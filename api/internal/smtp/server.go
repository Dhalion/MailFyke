package smtp

import (
	"fmt"
	"os"

	"github.com/Dhalion/MailFyke/internal/config"
	"github.com/Dhalion/MailFyke/internal/database/queries"
	"github.com/emersion/go-smtp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Server struct {
	srv    *smtp.Server
	pool   *pgxpool.Pool
	q      *queries.Queries
	config *config.Config
	log    zerolog.Logger
}

func New(pool *pgxpool.Pool, cfg *config.Config) *Server {
	q := queries.New(pool)
	be := &Backend{
		pool:   pool,
		q:      q,
		config: cfg,
	}

	return &Server{
		log:    zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger(),
		srv:    smtp.NewServer(be),
		pool:   pool,
		q:      queries.New(pool),
		config: cfg,
	}
}

func (s *Server) Start(addr string) error {
	s.srv.AllowInsecureAuth = s.config.SMTPAllowInsecure
	s.srv.Addr = addr
	s.srv.Domain = s.config.SMTPDomain

	if s.config.SMTPDebug {
		s.srv.Debug = zerolog.ConsoleWriter{Out: os.Stdout}
	}

	go func() {
		s.log.Info().Msg(fmt.Sprintf("Starting SMTP Server on %s:%s", s.config.SMTPDomain, s.config.SMTPListenAddr))
		s.log.Info().Msg(fmt.Sprintf("Settings: %+v", s.config))
		err := s.srv.ListenAndServe()
		if err != nil {
			s.log.Error().Err(err).Msg("smtp server error")
			panic(err)
		}
		s.log.Info().Msg("smtp server stopped")
	}()
	return nil
}

func (s *Server) Stop() error {
	err := s.srv.Close()
	if err != nil {
		return err
	}
	return nil
}
