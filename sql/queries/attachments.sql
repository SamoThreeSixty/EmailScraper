-- name: SaveAttachment :one
INSERT INTO attachments (email_id, type, filename, path)
VALUES ($1, $2, $3, $4)
RETURNING *;