package main

import (
	"fmt"

	"github.com/gclinoz/Gator-go/internal/config"
)

type state struct {
	ptr	*config.Config
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
	if len(cmd.args) == 0 {
		return fmt.Errorf("expect a single argument, none provided")
	}

	err := s.ptr.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Println("login", cmd.args[0], "success!")

	return nil
}
