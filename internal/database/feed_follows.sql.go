// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: feed_follows.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeedFollow = `-- name: CreateFeedFollow :one
WITH ff AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES($1, $2, $3, $4, $5)
    RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id,
    f.name AS feed_name,
    u.name AS user_name
FROM ff,
    feeds f,
    users u
WHERE ff.user_id = u.id
    AND ff.feed_id = f.id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
	FeedName  string
	UserName  string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (CreateFeedFollowRow, error) {
	row := q.db.QueryRowContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i CreateFeedFollowRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
		&i.FeedName,
		&i.UserName,
	)
	return i, err
}

const deleteFollowFeed = `-- name: DeleteFollowFeed :exec
DELETE FROM feed_follows ff
WHERE ff.user_id = $1
    AND ff.feed_id = $2
`

type DeleteFollowFeedParams struct {
	UserID uuid.UUID
	FeedID uuid.UUID
}

func (q *Queries) DeleteFollowFeed(ctx context.Context, arg DeleteFollowFeedParams) error {
	_, err := q.db.ExecContext(ctx, deleteFollowFeed, arg.UserID, arg.FeedID)
	return err
}

const getFeedFollowsByUser = `-- name: GetFeedFollowsByUser :many
SELECT f.id, f.created_at, f.updated_at, f.name, f.url, f.user_id, f.last_fetched_at,
    u.name AS created_by
FROM feeds f,
    feed_follows ff,
    users u
WHERE f.id = ff.feed_id
    AND f.user_id = u.id
    AND ff.user_id = $1
`

type GetFeedFollowsByUserRow struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	Url           string
	UserID        uuid.UUID
	LastFetchedAt sql.NullTime
	CreatedBy     string
}

func (q *Queries) GetFeedFollowsByUser(ctx context.Context, userID uuid.UUID) ([]GetFeedFollowsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedFollowsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedFollowsByUserRow
	for rows.Next() {
		var i GetFeedFollowsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
			&i.LastFetchedAt,
			&i.CreatedBy,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
