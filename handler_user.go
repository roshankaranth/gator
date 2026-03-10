package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/roshankaranth/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("No username provided!\n")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return fmt.Errorf("User does not exist!\n")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("user %s has been set!\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("No username provided!\n")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err == nil {
		return fmt.Errorf("User already exists!")
	}

	user := database.CreatedUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
	}

	s.db.CreatedUser(context.Background(), user)
	s.cfg.SetUser(cmd.args[0])
	fmt.Printf("%s user created!\n", cmd.args[0])

	return nil
}
