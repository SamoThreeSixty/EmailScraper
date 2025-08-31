-- +goose Up

ALTER TABLE email RENAME COLUMN body TO html_body;
ALTER TABLE email ADD COLUMN text_body TEXT;
ALTER TABLE email DROP COLUMN created_at;
ALTER TABLE email ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW();

-- +goose Down

ALTER TABLE email RENAME COLUMN html_body TO body;
ALTER TABLE email DROP COLUMN text_body;
