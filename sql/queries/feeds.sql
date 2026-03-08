-- name: Createfeed :one
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

-- name: Getfeed :one
SELECT *
FROM feeds
WHERE name = $1
LIMIT 1;    

-- name: Deletefeeds :exec
DELETE FROM feeds;      

-- name: GetAllfeeds :many
SELECT *
FROM feeds;
 
-- name: Getfeedbyurl :one
SELECT *
FROM feeds
WHERE url = $1
LIMIT 1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;