-- name: GetRestaurantUserRoleID :one
SELECT role_id FROM restaurant_users 
WHERE restaurant_id = $1 AND user_id = $2;

-- name: GetAnyRestaurantUserLinkByUserID :one
SELECT * FROM restaurant_users
WHERE user_id = $1
LIMIT 1;

-- name: AddRestaurantUser :exec
INSERT INTO restaurant_users (user_id, restaurant_id, role_id)
VALUES ($1, $2, $3);
