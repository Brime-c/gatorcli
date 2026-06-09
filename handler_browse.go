package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Brime/gatorcli/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error{
	limit := 2

	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		limit = parsedLimit
	}
	posts, err := s.db.GetPostsForUser(context.Background(),database.GetPostsForUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := posts {
		fmt.Println("Title:", post.Title)
		fmt.Println("URL:", post.Url)
		fmt.Println("Published:", post.PublishedAt)
		fmt.Println()
	}
	
	return nil
}