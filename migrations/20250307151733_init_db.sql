-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users(
    id uuid DEFAULT gen_random_uuid(),
    telegram_id BIGINT,
    username VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255),
    registered_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    UNIQUE(telegram_id, username)
);

CREATE TABLE IF NOT EXISTS api_clients(
    id uuid DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    token uuid DEFAULT gen_random_uuid(),
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS bind_tokens(
    token uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
            REFERENCES users(id)
                ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS api_clients;
DROP TABLE IF EXISTS bind_tokens;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd


