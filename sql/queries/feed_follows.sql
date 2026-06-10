-- name: CreateFeedFollow :one
-- CreateFeedFollow inserts a new follow relationship between a user and a feed.
-- It uses a Common Table Expression (CTE) to write the follow record first,
-- then joins the feeds and users tables to return human-readable names alongside the IDs.
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES(
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
    feeds.name as feed_name,
    users.name as user_name
from inserted_feed_follow
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
-- GetFeedFollowsForUser retrieves all feed follow records for a specific user,
-- joining on the feeds and users tables to include names in the final results.
SELECT
    feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
from feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedByUserFeed :exec
-- DeleteFeedByUserFeed removes a specific feed follow record.
-- This effectively unsubscribes a user from the target feed.
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;