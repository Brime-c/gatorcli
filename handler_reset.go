package main

import (
	"context"
	"fmt"
)

// handlerReset purges the database, removing all registered users and feeds.
// This is a destructive operation typically used during development.
func handlerReset(s *state, cmd command) error {
	// 1. Delete all users from the database.
	// Note: Depending on database constraints (ON DELETE CASCADE), this may
	// automatically clean up associated user data like feed follows.
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	// 2. Delete all tracked feeds from the database.
	err = s.db.DeleteFeeds(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("database has been reset")
	return nil
}
