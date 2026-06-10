package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Brime/gatorcli/internal/database"
	"github.com/google/uuid"
)

// handlerFollow subscribes the currently logged-in user to a feed.
// It looks up the feed by its URL first to retrieve its database ID.
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no command provided")
	}

	// 1. Find the feed by URL to get its database record
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	// 2. Create the feed follow relationship in the database
	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	// 3. Print confirmation details (feed name and following user)
	fmt.Println(follow.FeedName)
	fmt.Println(follow.UserName)
	return nil
}

// handlerFollowing prints a list of all feeds that the currently logged-in user follows.
func handlerFollowing(s *state, cmd command, user database.User) error {
	// 1. Retrieve the list of feed follows for this user from the database
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	// 2. Print the name of each followed feed
	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}
	return nil
}

// handlerUnfollow removes the feed subscription for the logged-in user.
// It resolves the feed URL to an ID and deletes the matching join table record.
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("No url provided")
	}

	// 1. Find the feed by URL to get its ID
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	// 2. Delete the record matching both this user's ID and the feed's ID
	err = s.db.DeleteFeedByUserFeed(context.Background(), database.DeleteFeedByUserFeedParams{
		FeedID: feed.ID,
		UserID: user.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
