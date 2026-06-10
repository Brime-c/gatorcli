package main

import (
	"context"

	"github.com/Brime/gatorcli/internal/database"
)

// middlewareLoggedIn is a decorator that wraps a command handler requiring an authenticated user.
// It intercepts the execution flow to fetch the current user from the database using the name
// stored in the local configuration. If the user exists, it passes them directly to the handler;
// otherwise, it intercepts and returns an error, preventing unauthorized execution.
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// 1. Look up the currently logged-in user in the database
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		// 2. Pass the validated user into the wrapped handler
		return handler(s, cmd, user)
	}
}
