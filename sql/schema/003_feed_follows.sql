-- +goose Up

-- feed_follows acts as a join table to represent the many-to-many relationship
-- between users and the RSS feeds they choose to follow.
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    -- fk_user_id links the follow record to a valid user. 
    -- If a user account is deleted, their follow records are cascade-deleted.
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    feed_id UUID NOT NULL,
    -- fk_feed_id links the follow record to a valid feed. 
    -- If a feed is deleted, all follow records for that feed are cascade-deleted.
    CONSTRAINT fk_feed_id
        FOREIGN KEY (feed_id)
        REFERENCES feeds(id)
        ON DELETE CASCADE,
    -- This composite unique constraint ensures a user cannot follow the same feed multiple times.
    UNIQUE(user_id, feed_id)    
);

-- +goose Down

-- Safely tear down the join table during rollbacks.
DROP TABLE feed_follows;