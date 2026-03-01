-- name: ListProjectsByOwnerWithRoot :many
SELECT sqlc.embed(projects), sqlc.embed(topics) FROM projects
JOIN topics ON topics.project_id = projects.id AND topics.parent_id IS NULL
WHERE owner_id = $1;

-- name: GetProjectByID :one
SELECT sqlc.embed(projects), sqlc.embed(topics), sqlc.embed(users) FROM projects
JOIN topics ON topics.project_id = projects.id AND topics.parent_id IS NULL
JOIN users ON projects.owner_id = users.id
WHERE projects.id = $1;

-- name: CreateProject :one
INSERT INTO projects (owner_id) VALUES($1)
RETURNING *;

-- name: RemoveProjectByID :exec
DELETE FROM projects WHERE id = $1;