package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	ctx := context.Background()
	_, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil{
		return fmt.Errorf("can't login to an account that doesn't exist!")
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("login", cmd.args[0], "success!")
	return nil
}

func handlerRegis(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	ctx := context.Background()
	userInfo := database.CreateUserParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cmd.args[0],
	}
	_, err := s.db.CreateUser(ctx, userInfo)
	if err != nil {
		return fmt.Errorf("username already exists")
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println(cmd.args[0], "created!")
	data, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil{
		return err
	}
	fmt.Println(data)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		fmt.Println("additional arguments will be ignored")
	}

	ctx := context.Background()
	err := s.db.DeleteAllUser(ctx)
	if err != nil{
		return err
	}
	return nil
}
