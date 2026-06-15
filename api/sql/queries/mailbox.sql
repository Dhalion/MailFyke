-- name: ListEmails :many
SELECT * FROM emails
WHERE organization_id = $1
ORDER BY received_at DESC
LIMIT $2 OFFSET $3;

-- name: GetEmail :one
SELECT * FROM emails WHERE id = $1 AND organization_id = $2 LIMIT 1;

-- name: DeleteEmail :exec
DELETE FROM emails WHERE id = $1 AND organization_id = $2;

-- name: MarkEmailRead :exec
UPDATE emails SET read = $3 WHERE id = $1 AND organization_id = $2;

-- name: UnreadCount :one
SELECT COUNT(*) FROM emails WHERE organization_id = $1 AND read = false;

-- name: InsertEmail :one
INSERT INTO emails (organization_id, smtp_username, mail_from, rcpt_to, subject, body_html, body_text, headers_json, raw_eml, has_attachments, size_bytes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;
