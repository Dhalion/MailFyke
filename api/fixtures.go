package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Dhalion/MailFyke/internal/config"
	"github.com/Dhalion/MailFyke/internal/database"
	"github.com/Dhalion/MailFyke/internal/database/queries"
	"golang.org/x/crypto/bcrypt"
)

func runFixtures(cfg *config.Config) error {
	ctx := context.Background()

	pool, err := database.NewPool(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer pool.Close()

	q := queries.New(pool)

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash: %w", err)
	}

	org, err := q.CreateOrganization(ctx, queries.CreateOrganizationParams{
		Name: "Acme Corp",
		Slug: "acme",
	})
	if err != nil {
		return fmt.Errorf("create org: %w", err)
	}
	log.Printf("created organization %s (%s)", org.Name, org.Slug)

	admin, err := q.CreateUser(ctx, queries.CreateUserParams{
		Email:        "admin@mailfyke.dev",
		PasswordHash: string(hash),
		IsAdmin:      true,
	})
	if err != nil {
		return fmt.Errorf("create admin: %w", err)
	}
	log.Printf("created admin user %s", admin.Email)

	if err := q.AddUserToOrganization(ctx, queries.AddUserToOrganizationParams{
		UserID:         admin.ID,
		OrganizationID: org.ID,
		Role:           "member",
	}); err != nil {
		return fmt.Errorf("add admin to org: %w", err)
	}
	log.Printf("added admin to organization")

	smtpHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("smtp hash: %w", err)
	}

	cred, err := q.CreateSMTPCredential(ctx, queries.CreateSMTPCredentialParams{
		OrganizationID: org.ID,
		Username:       "tenant_a",
		PasswordHash:   string(smtpHash),
	})
	if err != nil {
		return fmt.Errorf("create smtp cred: %w", err)
	}
	log.Printf("created SMTP credential %s", cred.Username)

	userHash, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("user hash: %w", err)
	}

	user, err := q.CreateUser(ctx, queries.CreateUserParams{
		Email:        "user@mailfyke.dev",
		PasswordHash: string(userHash),
		IsAdmin:      false,
	})
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	log.Printf("created regular user %s", user.Email)

	if err := q.AddUserToOrganization(ctx, queries.AddUserToOrganizationParams{
		UserID:         user.ID,
		OrganizationID: org.ID,
		Role:           "member",
	}); err != nil {
		return fmt.Errorf("add user to org: %w", err)
	}
	log.Printf("added user to organization")

	sampleEmails := []struct {
		from    string
		rcpt    string
		subject string
		body    string
	}{
		{"alice@example.com", "info@acme.com", "Welcome to Acme Corp", "Hi there, welcome aboard!"},
		{"bob@external.com", "support@acme.com", "Support Request #42", "I need help with my account."},
		{"newsletter@news.com", "info@acme.com", "Monthly Newsletter", "Here is your monthly update."},
	}

	for _, e := range sampleEmails {
		headers := fmt.Sprintf(`{"From":"%s","To":"%s","Subject":"%s"}`, e.from, e.rcpt, e.subject)
		raw := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", e.from, e.rcpt, e.subject, e.body)

		if _, err := q.InsertEmail(ctx, queries.InsertEmailParams{
			OrganizationID: org.ID,
			SmtpUsername:   cred.Username,
			MailFrom:       e.from,
			RcptTo:         []string{e.rcpt},
			Subject:        e.subject,
			BodyHtml:       "",
			BodyText:       e.body,
			HeadersJson:    []byte(headers),
			RawEml:         raw,
			HasAttachments: false,
			SizeBytes:      int(len(raw)),
		}); err != nil {
			return fmt.Errorf("insert email: %w", err)
		}
	}
	log.Printf("inserted %d sample emails", len(sampleEmails))

	log.Printf("fixtures complete")
	log.Printf("  admin: admin@mailfyke.dev / admin123")
	log.Printf("  user:  user@mailfyke.dev / user123")
	log.Printf("  smtp:  tenant_a / secret123")
	return nil
}
