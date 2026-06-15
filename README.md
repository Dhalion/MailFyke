# MailFyke

Multi-tenant mail catcher. Receive emails via SMTP with tenant isolation,
view them in a web UI with role-based access.

## Quick Start

```bash
mise run setup          # install dependencies
mise run up             # start PostgreSQL
mise run migrate        # run migrations
mise run fixtures       # seed test data
mise run dev            # start API with hot-reload
```

Send a test email:

```bash
swaks --to test@example.com \
      --server localhost:2525 \
      --auth LOGIN \
      --auth-user <tenant_username> \
      --auth-password <tenant_password> \
      --header "Subject: Test" \
      --body "Hello from MailFyke!"
```

