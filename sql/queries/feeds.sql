-- name: CreateFeed :one
-- CreateFeed inserts a new RSS feed into the database.
-- It returns the complete row of the newly registered feed.
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

-- name: DeleteFeeds :exec
-- DeleteFeeds purges all feeds from the database.
-- This is a destructive operation used by the developer 'reset' command.
DELETE FROM feeds;

-- name: ListFeedsWithName :many
-- ListFeedsWithName retrieves all registered feeds, joining them with the users
-- table to include the username of the user who registered each feed.
SELECT  feeds.name AS feed_name,
feeds.url,
users.name AS user_name
FROM feeds
JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
-- GetFeedByURL retrieves a single feed record by its source URL.
-- This is primarily used to resolve a URL to a feed ID during a 'follow' or 'unfollow' action.
SELECT * FROM feeds
WHERE URL = $1;

-- name: MarkFeedFetched :exec
-- MarkFeedFetched updates a feed's last_fetched_at and updated_at timestamps to the current time.
-- This signals to the scheduler that this feed has just been successfully crawled.
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
-- GetNextFeedToFetch selects the single feed that is most due to be scraped.
-- It returns the feed with the oldest 'last_fetched_at' timestamp,
-- prioritizing feeds that have never been fetched (which have NULL values).
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;