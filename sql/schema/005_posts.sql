-- +goose Up

-- The posts table stores individual articles, blogs, or episodes 
-- parsed and scraped from the RSS feeds.
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    -- The URL must be unique to prevent storing or displaying duplicate posts
    -- from the same RSS feed during consecutive scrape cycles.
    url TEXT UNIQUE NOT NULL,
    description TEXT,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    -- fk_feed_id ensures referential integrity; posts cannot exist without a valid parent feed.
    -- ON DELETE CASCADE ensures that if a feed is removed, all of its scraped posts are cleaned up too.
    CONSTRAINT fk_feed_id
        FOREIGN KEY (feed_id)
        REFERENCES feeds(id)
        ON DELETE CASCADE
);

-- +goose Down

-- Safely remove the posts table during a rollback.
DROP TABLE posts;