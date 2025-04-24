-- +goose Up
-- +goose StatementBegin
CREATE TABLE restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    name VARCHAR(50) NOT NULL,
    alias VARCHAR(50) NOT NULL,
    description TEXT NULL,
    address VARCHAR(255) NOT NULL,
    lat FLOAT NULL,
    lng FLOAT NULL,
    phone VARCHAR(30) NULL,
    image_url VARCHAR(255) NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    place_id VARCHAR(50) NOT NULL
);

CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON restaurants
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_updated_at ON restaurants;
DROP TABLE IF EXISTS restaurants;
-- +goose StatementEnd
