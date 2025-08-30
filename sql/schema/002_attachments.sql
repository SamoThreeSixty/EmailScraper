-- +goose UP

CREATE TABLE attachments
(
    id       SERIAL PRIMARY KEY NOT NULL,
    email_id INTEGER REFERENCES email (id) NOT NULL,
    type     VARCHAR(255) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    path     VARCHAR(255) NOT NULL
);

-- +goose DOWN

DROP TABLE attachments;