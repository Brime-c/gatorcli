-- +goose Up

-- The feeds table stores the metadata and source URLs of the RSS feeds
-- registered in our system by users.
CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    -- The source URL of the RSS feed. This is marked UNIQUE to ensure
    -- that we only store and poll a single instance of each feed resource.
    url TEXT UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    -- fk_user_id references the user who originally added/registered this feed.
    -- If the user's account is deleted, the feeds they registered are cascadingly removed.
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down

-- Safely drop the feeds table during a rollback.
DROP TABLE feeds;