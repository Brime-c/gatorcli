package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Brime/gatorcli/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// handlerAgg starts an infinite background loop that periodically scrapes RSS feeds.
// It parses a duration string (e.g. "1m", "1h") from CLI arguments to define the interval.
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no time between provided")
	}

	// 1. Parse duration string into a Go time.Duration
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	// 2. Start a ticker to run the scrape task on the defined interval
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println(err)
		}
	}
}

// scrapeFeeds identifies the next feed that is due to be fetched, marks it as fetched,
// retrieves its XML content, and persists new feed items as posts in the database.
func scrapeFeeds(s *state) error {
	ctx := context.Background()

	// 1. Fetch the oldest/least-recently-fetched feed from the database queue
	next, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	// 2. Instantly update its timestamps so other workers/processes don't scrape it concurrently
	err = s.db.MarkFeedFetched(ctx, next.ID)
	if err != nil {
		return err
	}

	// 3. Make HTTP request and parse the XML RSS feed
	feed, err := fetchFeed(ctx, next.Url)
	if err != nil {
		return err
	}

	// 4. Iterate over items in the feed and save them to the database
	for _, item := range feed.Channel.Item {
		publishedAt, err := parseTime(item.PubDate)
		if err != nil {
			log.Println("error parsing published time:", err)
			continue
		}

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: publishedAt,
			FeedID:      next.ID,
		})
		if err != nil {
			// Check if the error is a PostgreSQL unique constraint violation (duplicate URL).
			// If it is, we silently skip saving this post as it has already been scraped.
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			log.Println("error creating post:", err)
		}
	}
	return nil
}

// parseTime attempts to parse an RSS pubDate string using a list of common feed date formats.
func parseTime(s string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse time: %q", s)
}
