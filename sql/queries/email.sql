-- name: GetEmail :one
SELECT *
FROM email
WHERE id = $1;

-- name: InsertEmail :exec
INSERT INTO email (id, subjectfrom_email, to_email, date_sent, body, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;