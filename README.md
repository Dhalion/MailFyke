<p align="center">
  <img src="docs/logo.png" alt="MailFyke" width="250">
</p>

<h1 align="center">MailFyke</h1>

<p align="center">
  Multi-tenant mail catcher — SMTP receiver with tenant isolation and web UI
  <br>
  <strong>Go</strong> · <strong>Chi</strong> · <strong>PostgreSQL</strong> · <strong>Vue 3</strong> · <strong>Tailwind</strong>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-blue?logo=go" alt="Go 1.26">
  <img src="https://img.shields.io/badge/Node-26-green?logo=node.js" alt="Node 26">
  <img src="https://img.shields.io/badge/status-alpha-orange" alt="Status Alpha">
</p>

---

## Features

- **SMTP receiver** — catch-all inbound emails on port `:2525` (AUTH LOGIN/PLAIN)
- **Multi-tenant** — tenant derived from SMTP credentials, complete mailbox isolation
- **Web UI** — browse, read, and manage emails per mailbox
- **Role-based access** — organization-scoped users with admin/member roles
- **Single binary** — SMTP + HTTP in one process, embedded frontend
- **Local-first** — PostgreSQL via Docker Compose, no cloud dependencies

## Quick Start

```bash
mise run setup      # install Go + Node dependencies
mise run up         # start PostgreSQL
mise run migrate    # run database migrations
mise run fixtures   # seed test data
mise run dev        # start API with hot-reload
```

## Send a test email

```bash
swaks --to test@example.com \
      --server localhost:2525 \
      --auth PLAIN \
      --auth-user tenant_a \
      --auth-password secret123 \
      --header "Subject: Hello" \
      --body "Welcome to MailFyke!"
```

Open [http://localhost:5789](http://localhost:5789) and log in with `admin@mailfyke.dev` / `password`.

## Architecture

```
SMTP client ──► :2525 ──► go-smtp ──► Auth (bcrypt) ──► PostgreSQL
Browser     ──► :5789 ──► chi HTTP ──► JWT middleware ──► PostgreSQL
```

One process, one database pool. SMTP auth looks up `smtp_credentials`, derives `organization_id`. HTTP API validates JWT, enforces org membership. Both share generated sqlc queries.

## Configuration

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://mailfyke:mailfyke@localhost:5432/mailfyke?sslmode=disable` | PostgreSQL connection |
| `LISTEN_ADDR` | `:5789` | HTTP listen address |
| `SMTP_LISTEN_ADDR` | `:2525` | SMTP listen address |
| `SMTP_ALLOW_INSECURE` | `true` | Allow AUTH without TLS |
| `SMTP_DOMAIN` | `localhost` | SMTP banner domain |
| `SMTP_DEBUG` | `false` | Log SMTP protocol to stdout |
| `JWT_SECRET` | `dev-secret-change-in-production` | HMAC key for JWT tokens |

## Development

```bash
mise run generate   # Regenerate sqlc + oapi-codegen + frontend types
mise run lint       # Run golangci-lint + oxlint
mise run test       # Run Go tests
mise run build      # Build binary to tmp/mailfyke
```

## Project Structure

```
api/                      # Go backend
├── internal/
│   ├── smtp/             # SMTP receiver (go-smtp)
│   ├── handler/          # HTTP handlers (oapi-codegen)
│   ├── middleware/        # JWT auth, org membership
│   ├── database/queries/ # sqlc generated — DO NOT EDIT
│   ├── config/           # Env-based config
│   └── mail/             # Email parsing (enmime)
├── sql/schema.sql        # DDL (source of truth)
├── sql/queries/          # sqlc query definitions
├── openapi/spec.yaml     # OpenAPI 3.0 spec
└── server.go             # Router wiring + server lifecycle

frontend/                 # Vue 3 + Tailwind + Reka UI
docs/logo.png             # Project logo
compose.yml               # PostgreSQL service
Dockerfile                # Distroless container build
```

## Testing SMTP

Manual dialog:

```bash
telnet localhost 2525
EHLO client
AUTH LOGIN
# base64(username)
# base64(password)
MAIL FROM:<test@example.com>
RCPT TO:<user@example.com>
DATA
Subject: Test
From: test@example.com
To: user@example.com

Hello from MailFyke!
.
QUIT
```