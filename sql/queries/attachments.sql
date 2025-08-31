-- name: SaveAttachment :one
INSERT INTO attachments (email_id, created_at, type, original_filename, saved_filename, path)
VALUES ($1, NOW(), $2, $3, $4, $5)
RETURNING *;

-- name: GetEmailAttachments :many
SELECT * FROM attachments WHERE email_id = $1;

-- name: UpdateAttachmentPathFilenames :exec
UPDATE attachments SET path = $1, saved_filename = $2, original_filename = $3 WHERE id = $4;