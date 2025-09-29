-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_verifications (
    user_id INT NOT NULL PRIMARY KEY,
    code VARCHAR(6) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_verifications;
-- +goose StatementEnd
