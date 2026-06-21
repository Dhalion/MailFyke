-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, is_admin) VALUES ($1, $2, $3) RETURNING *;

-- name: GetAccountByUsername :one
SELECT a.*, o.slug AS organization_slug
FROM smtp_accounts a
JOIN organizations o ON o.id = a.organization_id
WHERE a.username = $1 LIMIT 1;

-- name: CreateAccount :one
INSERT INTO smtp_accounts (organization_id, username)
VALUES ($1, $2) RETURNING *;

-- name: CreateAuthMethod :one
INSERT INTO smtp_auth_methods (account_id, type, hash, description, expires_at)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetAuthMethodsByAccountID :many
SELECT * FROM smtp_auth_methods WHERE account_id = $1;

-- name: GetPasswordMethodByAccountID :one
SELECT * FROM smtp_auth_methods WHERE account_id = $1 AND type = 'password' LIMIT 1;

-- name: DeleteAuthMethod :exec
DELETE FROM smtp_auth_methods WHERE id = $1;

-- name: AddUserToOrganization :exec
INSERT INTO user_organization (organization_id, user_id, role) VALUES ($1, $2, $3);

-- name: RemoveUserFromOrganization :exec
DELETE FROM user_organization WHERE organization_id = $1 AND user_id = $2;

