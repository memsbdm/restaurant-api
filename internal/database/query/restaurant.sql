-- name: IsRestaurantAlreadyTaken :one
SELECT EXISTS (
    SELECT 1
    FROM restaurants r
    LEFT JOIN restaurant_users ru ON ru.restaurant_id = r.id AND ru.user_id = $2
    WHERE r.place_id = $1
    AND (
        r.is_verified = TRUE OR ru.id IS NOT NULL
    )
);

-- name: CreateRestaurant :one
INSERT INTO restaurants
(name, alias, address, lat, lng, phone, place_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
