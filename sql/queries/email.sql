-- name: GetEmail :one
SELECT *
FROM email
WHERE id = $1;

-- name: InsertEmail :one
INSERT INTO email (created_at, subject, from_email, to_email, date_sent, html_body, text_body)
VALUES (NOW(), $1, $2, $3, $4, $5, $6)
RETURNING *;