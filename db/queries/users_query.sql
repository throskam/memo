-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByAuthenticationMethod :one
SELECT sqlc.embed(users), sqlc.embed(authentication_methods) FROM users
JOIN authentication_methods ON authentication_methods.user_id = users.id
WHERE authentication_methods.provider = $1 AND authentication_methods.sub = $2
LIMIT 1;

-- name: CreateUser :one

INSERT INTO users (username) VALUES($1)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2, updated_at = now()
WHERE id = $1
RETURNING *;