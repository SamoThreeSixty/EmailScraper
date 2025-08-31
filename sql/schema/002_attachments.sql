-- +goose UP

CREATE TABLE attachments
(
    id                SERIAL PRIMARY KEY            NOT NULL,
    email_id          INTEGER REFERENCES email (id) NOT NULL,
    created_at        TIMESTAMP                     NOT NULL DEFAULT NOW(),
    type              VARCHAR(255)                  NOT NULL,
    original_filename VARCHAR(255)                  NOT NULL,
    saved_filename    VARCHAR(255)                  NOT NULL,
    path              VARCHAR(255)                  NOT NULL
);

-- +goose DOWN

DROP TABLE attachments;