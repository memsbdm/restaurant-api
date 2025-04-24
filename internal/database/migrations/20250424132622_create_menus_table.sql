-- +goose Up
-- +goose StatementBegin
CREATE TABLE menus (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  menu_order SMALLINT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS menus;
-- +goose StatementEnd
