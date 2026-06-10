package main

import (
	"log"
	"context"
	"github.com/gclinoz/Gator-go/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error { 
		user, err := s.db.GetUser(context.Background(), s.cfg.Username)
		if err != nil {
			log.Fatal("user does not exist")
		}
		handler(s, cmd, user)
		return nil
	}
}
