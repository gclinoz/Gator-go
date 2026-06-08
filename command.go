package main

import (
	"fmt"

	"github.com/gclinoz/Gator-go/internal/config"
	"github.com/gclinoz/Gator-go/internal/database"
)

type state struct {
	db	*database.Queries
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
