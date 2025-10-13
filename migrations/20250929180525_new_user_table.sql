-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING' NOT NULL,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX users_login_idx ON users (login);
CREATE INDEX users_email_idx ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX users_email_idx;
DROP INDEX users_login_idx;
DROP TABLE users;
-- +goose StatementEnd
