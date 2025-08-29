-- +goose UP

CREATE TABLE email (
    id SERIAL PRIMARY KEY,
    subject TEXT NOT NULL,
    from_email TEXT NOT NULL,
    to_email TEXT NOT NULL,
    date_sent TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose DOWN

DROP TABLE email;