-- name: CreatePost :one
-- CreatePost inserts a newly scraped blog post or article into the database.
-- It returns the fully populated post record.
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;


-- name: GetPostsForUser :many
-- GetPostsForUser retrieves a chronological timeline of posts for a specific user.
-- It joins the posts and feed_follows tables on feed_id, filtering for feeds 
-- followed by the target user, and orders the results by publication date (newest first).
SELECT posts.* 
FROM posts
JOIN feed_follows on posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;