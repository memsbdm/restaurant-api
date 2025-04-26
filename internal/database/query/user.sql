-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET email = $1, is_email_verified = $2, avatar_url= $3
WHERE id = $4
RETURNING *;

-- name: UserEmailTaken :one
SELECT EXISTS(
  SELECT 1 FROM users
  WHERE email = $1 AND is_email_verified = TRUE
);

