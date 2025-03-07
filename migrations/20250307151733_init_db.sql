-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS users(
                                    id serial PRIMARY KEY,
                                    name VARCHAR(255),
                                    password VARCHAR(255),
                                    registered_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    UNIQUE(name)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd


