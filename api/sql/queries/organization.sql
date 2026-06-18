-- name: ListOrganizationsByUser :many
SELECT o.*, uo.role FROM organizations o
JOIN user_organization uo ON uo.organization_id = o.id
WHERE uo.user_id = $1;

-- name: ListOrganizationsByAdmin :many
SELECT * FROM organizations ORDER BY name;

-- name: GetOrganizationBySlug :one
SELECT * FROM organizations WHERE slug = $1 LIMIT 1;

-- name: CreateOrganization :one
INSERT INTO organizations (name, slug) VALUES ($1, $2) RETURNING *;
