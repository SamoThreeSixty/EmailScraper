-- +goose UP

CREATE TABLE email
(
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    subject    TEXT      NOT NULL,
    from_email TEXT      NOT NULL,
    to_email   TEXT      NOT NULL,
    date_sent  TIMESTAMP NOT NULL,
    html_body  TEXT      NOT NULL,
    text_body  TEXT      NOT NULL
);

-- +goose DOWN

DROP TABLE email;