-- name: SaveAttachment :one
INSERT INTO attachments (email_id, created_at, type, original_filename, saved_filename, path)
VALUES ($1, NOW(), $2, $3, $4, $5)
RETURNING *;

-- name: GetEmailAttachments :many
SELECT * FROM attachments WHERE email_id = $1;