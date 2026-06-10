package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Brime/gatorcli/internal/database"
)

// handlerBrowse displays a paginated list of posts from feeds followed by the logged-in user.
// It accepts an optional argument to override the default limit of posts returned.
func handlerBrowse(s *state, cmd command, user database.User) error {
	// 1. Establish a default maximum number of posts to fetch
	limit := 2

	// 2. If the user provided an argument, parse it as an integer to override the default limit
	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = parsedLimit
	}

	// 3. Fetch the posts for the user's followed feeds up to the specified limit
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	// 4. Print each post to the terminal with clean formatting
	for _, post := range posts {
		fmt.Println("Title:", post.Title)
		fmt.Println("URL:", post.Url)
		fmt.Println("Published:", post.PublishedAt)
		fmt.Println()
	}

	return nil
}
