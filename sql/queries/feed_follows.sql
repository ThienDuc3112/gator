-- name: CreateFeedFollow :one
WITH ff AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT ff.*,
    f.name AS feed_name,
    u.name AS user_name
FROM ff,
    feeds f,
    users u
WHERE ff.user_id = u.id
    AND ff.feed_id = f.id;
-- name: GetFeedFollowsByUser :many
SELECT f.*,
    u.name AS created_by
FROM feeds f,
    feed_follows ff,
    users u
WHERE f.id = ff.feed_id
    AND f.user_id = u.id
    AND ff.user_id = $1;
-- name: DeleteFollowFeed :exec
DELETE FROM feed_follows ff
WHERE ff.user_id = $1
    AND ff.feed_id = $2;