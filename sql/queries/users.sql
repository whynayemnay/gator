-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE name = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT name FROM users;

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
SELECT feeds.name, feeds.url, users.name FROM feeds 
    INNER JOIN users
    ON feeds.user_id = users.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: SelectFeedByURL :one
SELECT * FROM feeds
WHERE feeds.url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,  -- Select all columns from feed_follows table
    feeds.name AS feed_name,  -- Select feed name from the feeds table
    users.name AS user_name  -- Select user name from the users table
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE feed_follows.user_id = $1;  -- Filter by user_id (replace $1 with the specific user's ID)

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE feed_id = $1 AND user_id = $2;
--
