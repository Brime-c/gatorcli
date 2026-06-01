package main

import (
	"fmt"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Errorf("no command provided")
	}
}
