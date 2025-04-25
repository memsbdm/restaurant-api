-- name: AddRestaurantUser :exec
INSERT INTO restaurant_users (user_id, restaurant_id, role_id)
VALUES ($1, $2, $3);
