-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, feed_id, title, url,
description, published_at) 
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
);

-- name: GetPostUser :many
SELECT posts.*, feeds.name
FROM posts
INNER JOIN feed_follows on posts.feed_id = feed_follows.feed_id
INNER JOIN feeds ON posts.feed_id = feeds.id
WHERE feeds.user_id = $1
ORDER BY posts.created_at DESC
LIMIT $2;