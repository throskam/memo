-- name: CreateAuthenticationMethod :one
INSERT INTO authentication_methods (provider, sub, user_id) VALUES($1, $2, $3)
RETURNING *;