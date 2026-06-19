package smtp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/Dhalion/MailFyke/internal/config"
	"github.com/Dhalion/MailFyke/internal/database/queries"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jhillyerd/enmime/v2"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	backend      *Backend
	auth         bool
	pool         *pgxpool.Pool
	config       *config.Config
	q            *queries.Queries
	msg          *message
	orgId        pgtype.UUID
	smtpUsername string
}

type message struct {
	From     string
	To       []string
	RcptOpts []*smtp.RcptOptions
	Data     []byte
	envelope *enmime.Envelope
	Opts     *smtp.MailOptions
}

func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain}
}

func (s *Session) Auth(mech string) (sasl.Server, error) {
	return sasl.NewPlainServer(func(identity, username, password string) error {
		if identity != "" && identity != username {
			return errors.New("invalid identity")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		smtpCreds, err := s.q.GetSMTPCredentialByUsername(ctx, username)
		if err != nil {
			return errors.New("invalid credentials")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(smtpCreds.PasswordHash), []byte(password)); err != nil {
			return errors.New("invalid credentials")
		}

		s.auth = true
		s.smtpUsername = username
		s.orgId = smtpCreds.OrganizationID
		return nil
	}), nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	s.Reset()
	s.msg.From = from
	s.msg.Opts = opts
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.msg.To = append(s.msg.To, to)
	s.msg.RcptOpts = append(s.msg.RcptOpts, opts)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	rawBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	s.msg.Data = rawBytes
	s.msg.envelope, err = enmime.ReadEnvelope(bytes.NewReader(rawBytes))
	if err != nil {
		return err
	}

	err = persistMailInDb(s)
	if err != nil {
		return errors.New("internal error")
	}

	return nil
}

func (s *Session) Reset() {
	s.msg = &message{}
}

func (s *Session) Logout() error {
	return nil
}

func persistMailInDb(s *Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.q.InsertEmail(ctx, queries.InsertEmailParams{
		OrganizationID: s.orgId,
		SmtpUsername:   s.smtpUsername,
		MailFrom:       s.msg.From,
		RcptTo:         s.msg.To,
		Subject:        s.msg.envelope.GetHeader("Subject"),
		BodyHtml:       s.msg.envelope.HTML,
		BodyText:       s.msg.envelope.Text,
		HeadersJson:    []byte("{}"),
		RawEml:         string(s.msg.Data),
		HasAttachments: len(s.msg.envelope.Attachments) > 0,
		SizeBytes:      len(s.msg.Data),
	})

	if err != nil {
		return err
	}

	return nil
}
