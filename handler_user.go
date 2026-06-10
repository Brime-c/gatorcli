package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Brime/gatorcli/internal/database"
	"github.com/google/uuid"
)

// handlerLogin switches the active user in the local configuration file.
// It verifies that the requested user exists in the database before switching.
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided")
	}

	// 1. Verify the user exists in the database
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	// 2. Persist the username to the local config file
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("username has been set to %s\n", user.Name)
	return nil
}

// handlerRegister creates a new user in the database with a unique ID and timestamps.
// Upon successful registration, it automatically logs the new user in.
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided for registering")
	}

	// 1. Write the new user to the database
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("username already exists")
	}

	// 2. Log in automatically as the newly registered user
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("username %s was registered\n", user.Name)
	fmt.Printf("%+v\n", user)
	return nil
}

// handlerUsers lists all registered users in the database,
// highlighting the currently active user logged into the local configuration.
func handlerUsers(s *state, cmd command) error {
	// 1. Retrieve all users from the database
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	// 2. Print each user, adding a marker (*) and "(current)" next to the logged-in user
	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}
