-- name: MenuExistsForRestaurantID :one
SELECT EXISTS ( 
SELECT 1
FROM menus
WHERE restaurant_id = $1
AND is_active = TRUE
LIMIT 1
);

-- name: CreateMenu :one
INSERT INTO  menus (name, is_active, restaurant_id)
VALUES ($1, $2, $3)
RETURNING *;
