-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
