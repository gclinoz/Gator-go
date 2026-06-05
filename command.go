package main

import (
	"fmt"

	"github.com/gclinoz/Gator-go/internal/config"
)

type state struct {
	cfg	*config.Config
}

type command struct {
	name	string
	args	[]string
}

type commands struct {
	utils	map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.utils[cmd.name]
	if !ok {
		return fmt.Errorf("command not exists")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.utils[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("login", cmd.args[0], "success!")
	return nil
}
