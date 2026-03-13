package main

import (
	"context"

	"github.com/roshankaranth/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		current_user := s.cfg.Current_user_name
		user, err := s.db.GetUser(context.Background(), current_user)

		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
