CREATE TABLE organizations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_admin      BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_organization (
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            TEXT NOT NULL DEFAULT 'member',
    PRIMARY KEY (organization_id, user_id)
);

CREATE TABLE smtp_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    username        TEXT NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE smtp_auth_methods (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id  UUID NOT NULL REFERENCES smtp_accounts(id) ON DELETE CASCADE,
    type        TEXT NOT NULL,
    hash        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    expires_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_one_password_per_account
    ON smtp_auth_methods (account_id)
    WHERE type = 'password';

CREATE TABLE emails (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id  UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    smtp_username    TEXT NOT NULL,
    mail_from        TEXT NOT NULL,
    rcpt_to          TEXT[] NOT NULL,
    subject          TEXT NOT NULL DEFAULT '',
    body_html        TEXT NOT NULL DEFAULT '',
    body_text        TEXT NOT NULL DEFAULT '',
    headers_json     JSONB NOT NULL DEFAULT '{}',
    raw_eml          TEXT NOT NULL,
    has_attachments  BOOLEAN NOT NULL DEFAULT false,
    read             BOOLEAN NOT NULL DEFAULT false,
    received_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    size_bytes       INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_emails_org_received ON emails(organization_id, received_at DESC);
CREATE INDEX idx_emails_org_unread ON emails(organization_id, read) WHERE read = false;
