-- +goose Up

-- The users table acts as the core identity table, storing accounts
-- for everyone interacting with the gator CLI.
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    -- The username used for logging in and identification.
    -- This is UNIQUE to prevent registration conflicts.
    name TEXT UNIQUE NOT NULL
);

-- +goose Down

-- Safely drop the users table during rollbacks.
-- Warning: Because other tables (feeds, feed_follows) reference this table,
-- running this migration rollback will require dropping those dependent tables first.
DROP TABLE users;