-- name: ListTopicAncestors :many
WITH RECURSIVE ancestors AS (
    -- Base case: select the topic itself
    SELECT *
    FROM topics
    WHERE topics.id = $1

    UNION ALL

    -- Recursive case: select the parent of the current topic
    SELECT t.*
    FROM topics t
    INNER JOIN ancestors a ON t.id = a.parent_id
)
SELECT *
FROM ancestors;

-- name: ListTopicDescendants :many
WITH RECURSIVE descendants AS (
    -- Base case: select the topic itself
    SELECT *
    FROM topics
    WHERE topics.id = $1

    UNION ALL

    -- Recursive case: select the children of the current topic
    SELECT t.*
    FROM topics t
    INNER JOIN descendants a ON t.parent_id = a.id
)
SELECT *
FROM descendants
ORDER BY sort_order ASC;

-- name: ListTopicChildren :many
SELECT * FROM topics
WHERE topics.parent_id = $1;

-- name: GetTopicByID :one
SELECT sqlc.embed(topics), sqlc.embed(projects) FROM topics
JOIN projects ON topics.project_id = projects.id
WHERE topics.id = $1;

-- name: CreateTopic :one
INSERT INTO topics (title, content, sort_order, parent_id, project_id) VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateTopic :one
UPDATE topics
SET title = $2, content = $3, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: RemoveTopicByID :exec
DELETE FROM topics WHERE id = $1;

-- name: RemoveTopicsByProjectID :exec
DELETE FROM topics WHERE project_id = $1;

-- name: MoveTopic :exec
UPDATE topics
SET parent_id = $1, sort_order = $2
WHERE id = $3;

-- name: ShiftTopics :exec
UPDATE topics
SET sort_order = sort_order + @amount
WHERE parent_id = $1 and sort_order >= @start;

-- name: ReindexTopics :exec
WITH reordered_topics AS (
    SELECT
        id,
        parent_id,
        ROW_NUMBER() OVER (PARTITION BY parent_id ORDER BY sort_order) AS new_sort_order
    FROM topics
    WHERE project_id = $1
)
UPDATE topics
SET sort_order = reordered_topics.new_sort_order
FROM reordered_topics
WHERE topics.id = reordered_topics.id AND topics.project_id = $1;

-- name: GetLastSortOrder :one
SELECT sort_order
FROM topics
WHERE parent_id = $1
ORDER BY sort_order DESC
LIMIT 1;