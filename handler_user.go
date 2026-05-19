package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Brime/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])

	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)

	if err != nil {
		return err
	}
	fmt.Printf("username has been set to %s\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided for registering")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})

	if err != nil {
		return fmt.Errorf("username already exists")
	}
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("username %s was registered", user.Name)
	fmt.Printf("%+v\n", user)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
		} else {
			fmt.Printf("* %s\n", user)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {

	const fixedFeedUrl = "https://www.wagslane.dev/index.xml"

	feed, err := fetchFeed(context.Background(), fixedFeedUrl)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", feed)
	return nil
}
func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]

	if !ok {
		return fmt.Errorf("command doesnt exist")
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}
