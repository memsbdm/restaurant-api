-- +goose Up
-- +goose StatementBegin
CREATE TABLE restaurant_invites (
  id SERIAL PRIMARY KEY,
  restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
  invited_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  canceled_by_user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL,
  email VARCHAR(255) NOT NULL,
  role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  accepted_at TIMESTAMP NULL,
  canceled_at TIMESTAMP NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TRIGGER set_updated_at
  BEFORE UPDATE ON restaurant_invites
  FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_restaurant_invites_restaurant_id ON restaurant_invites (restaurant_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at ON restaurant_invites;
DROP INDEX IF EXISTS idx_restaurant_invites_restaurant_id;
DROP TABLE IF EXISTS restaurant_invites;
-- +goose StatementEnd
