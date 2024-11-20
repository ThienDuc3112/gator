-- name: CreatePost :one
INSERT INTO posts (
        id,
        created_at,
        updated_at,
        title,
        url,
        description,
        published_at,
        feed_id
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
-- name: GetPostForUser :many
SELECT p.*
FROM posts p,
    feed_follows ff
WHERE ff.feed_id = p.feed_id
    AND ff.user_id = $1
ORDER BY p.published_at DESC NULLS LAST
LIMIT $2;