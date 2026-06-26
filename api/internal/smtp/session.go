package smtp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Dhalion/MailFyke/internal/config"
	"github.com/Dhalion/MailFyke/internal/database/queries"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/golang-jwt/jwt/v5"
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
	return []string{sasl.Plain, XOAUTH2}
}

func (s *Session) Auth(mech string) (sasl.Server, error) {
	switch mech {
	case sasl.Plain:
		return sasl.NewPlainServer(func(identity, username, password string) error {
			if identity != "" && identity != username {
				return smtp.ErrAuthFailed
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			account, err := s.q.GetAccountByUsername(ctx, username)
			if err != nil {
				return smtp.ErrAuthFailed
			}

			passMethod, err := s.q.GetPasswordMethodByAccountID(ctx, account.ID)
			if err != nil {
				return smtp.ErrAuthFailed
			}

			if err := bcrypt.CompareHashAndPassword([]byte(passMethod.Hash), []byte(password)); err != nil {
				return smtp.ErrAuthFailed
			}

			s.auth = true
			s.smtpUsername = username
			s.orgId = account.OrganizationID
			return nil
		}), nil
	case XOAUTH2:
		return NewXOAuth2Server(func(opts XOAuth2Options) error {
			token, err := jwt.Parse(opts.TokenBytes, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(s.config.JWTSecret), nil
			})
			if err != nil {
				return smtp.ErrAuthFailed
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				return smtp.ErrAuthFailed
			}

			username, _ := claims["username"].(string)
			orgId, _ := claims["org_id"].(string)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			account, err := s.q.GetAccountByUsername(ctx, username)
			if err != nil {
				return smtp.ErrAuthFailed
			}

			// Reject Auth on unexpected OrgId mismatch
			if account.OrganizationID.String() != orgId {
				return smtp.ErrAuthFailed
			}

			s.auth = true
			s.smtpUsername = account.Username
			s.orgId = account.OrganizationID

			return nil
		}), nil
	default:
		return nil, smtp.ErrAuthUnsupported
	}
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}

	s.Reset()
	s.msg.From = from
	s.msg.Opts = opts
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}

	s.msg.To = append(s.msg.To, to)
	s.msg.RcptOpts = append(s.msg.RcptOpts, opts)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}

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
		SizeBytes:      int32(len(s.msg.Data)),
	})

	if err != nil {
		return err
	}

	return nil
}
