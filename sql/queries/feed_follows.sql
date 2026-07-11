-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *)

SELECT
a.*,
b.name as user_name,
c.name as feed_name
FROM inserted_feed_follow a
JOIN users b ON a.user_id=b.id
JOIN feeds c ON a.feed_id=c.id;

-- name: GetFeedFollowsForUser :many
SELECT
a.name as user_name,
c.name as feed_name
FROM users a
JOIN feed_follows b ON a.id=b.user_id
JOIN feeds c ON b.feed_id=c.id
WHERE a.name = $1;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;