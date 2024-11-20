-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetFeeds :many
SELECT f.*,
    u.name AS username
FROM feeds f,
    users u
WHERE f.user_id = u.id;
-- name: GetFeedByUrl :one
SELECT f.*,
    u.name AS username
FROM feeds f,
    users u
WHERE f.user_id = u.id
    AND f.url = $1;
-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1,
    updated_at = $2
WHERE id = $3;
-- name: GetNextFeedToFetch :one
SELECT f.*
FROM feeds f
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;