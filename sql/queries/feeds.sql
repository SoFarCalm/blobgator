-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)

RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedName :one
SELECT id, name FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET UPDATED_AT = NOW(), LAST_FECTHED_AT = NOW()
WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fecthed_at ASC NULLS FIRST LIMIT 1;