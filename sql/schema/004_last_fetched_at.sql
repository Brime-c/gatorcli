-- +goose Up

-- Add last_fetched_at to track when the aggregator last completed a scrape of this feed.
-- This column is nullable because newly created feeds will have a NULL value,
-- indicating they have never been fetched and should be prioritized in the scraping queue.
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down

-- Remove last_fetched_at column, rolling back the feed scraping queue tracking.
ALTER TABLE feeds DROP COLUMN last_fetched_at;