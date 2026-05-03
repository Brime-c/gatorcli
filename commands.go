package main

import (
	"github.com/Brime/gatorcli/internal/config"
	"github.com/Brime/gatorcli/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}
