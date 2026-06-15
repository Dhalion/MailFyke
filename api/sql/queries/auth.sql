-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (email, password_hash, is_admin) VALUES ($1, $2, $3) RETURNING *;

-- name: GetSMTPCredentialByUsername :one
SELECT sc.*, o.slug AS organization_slug
FROM smtp_credentials sc
JOIN organizations o ON o.id = sc.organization_id
WHERE sc.username = $1 LIMIT 1;

-- name: CreateSMTPCredential :one
INSERT INTO smtp_credentials (organization_id, username, password_hash) VALUES ($1, $2, $3) RETURNING *;

-- name: AddUserToOrganization :exec
INSERT INTO user_organization (organization_id, user_id, role) VALUES ($1, $2, $3);

-- name: RemoveUserFromOrganization :exec
DELETE FROM user_organization WHERE organization_id = $1 AND user_id = $2;
