package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no username provided")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("username has been set to %s", cmd.args[0])
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
