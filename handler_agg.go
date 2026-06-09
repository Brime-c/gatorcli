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

func handlerAgg(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("no time between provided")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			fmt.Println(err)
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	next, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(ctx, next.ID)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(ctx, next.Url)
	if err != nil {
		return err
	}

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
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				continue
			}
			log.Println("error creating post:", err)
		}
	}
	return nil
}

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
