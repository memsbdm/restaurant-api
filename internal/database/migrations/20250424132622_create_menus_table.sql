-- +goose Up
-- +goose StatementBegin
CREATE TABLE menus (
  id SERIAL PRIMARY KEY,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  name VARCHAR(50) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE
);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON menus
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at ON menus;
DROP TABLE IF EXISTS menus;
-- +goose StatementEnd
