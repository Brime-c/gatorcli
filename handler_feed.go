package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Brime/gatorcli/internal/database"
	"github.com/google/uuid"
)

// handlerAddFeed registers a new RSS feed in the database and automatically
// creates a record for the current user to follow that feed.
// It requires an authenticated user passed via middleware.
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("please provide a name and url")
	}

	// 1. Insert the new feed details into the database
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	// 2. Automatically make the creator follow their newly created feed
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}

// handlerFeeds retrieves and prints all registered feeds from the database,
// along with the name of the user who registered each feed.
func handlerFeeds(s *state, cmd command) error {
	// 1. Query the database for feeds, joining with the users table to get creator names
	feeds, err := s.db.ListFeedsWithName(context.Background())
	if err != nil {
		return err
	}

	// 2. Print the information for each feed to the terminal
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
		fmt.Println(feed.Url)
		fmt.Println(feed.UserName)
	}
	return nil
}
