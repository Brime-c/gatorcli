package main

import (
	"fmt"

	"github.com/Brime/gatorcli/internal/config"
	"github.com/Brime/gatorcli/internal/database"
)

// state holds the shared application dependencies, such as the active
// database connection pool and the user configuration settings.
type state struct {
	db  *database.Queries
	cfg *config.Config
}

// command represents a single instruction passed by the user in the CLI,
// containing the action name and any arguments provided.
type command struct {
	name string
	args []string
}

// commands manages the mapping of command names to their handler functions.
type commands struct {
	handlers map[string]func(*state, command) error
}

// run executes a registered command handler by its name.
// It returns an error if the command is unregistered or if the handler fails.
func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("command doesnt exist")
	}

	return handler(s, cmd)
}

// register associates a command name with a handler function.
// This populates the internal handlers map for later execution.
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}
