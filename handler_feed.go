package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/roshankaranth/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command) error {

	if len(cmd.args) < 2 {
		return fmt.Errorf("insufficient arguments!\n")
	}

	current_user := s.cfg.Current_user_name
	user, err := s.db.GetUser(context.Background(), current_user)

	if err != nil {
		return err
	}

	userID := user.ID
	user_feed := database.CreatedFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    userID,
	}
	feed, err := s.db.CreatedFeed(context.Background(), user_feed)

	if err != nil {
		return err
	}

	fmt.Printf("Feed created succesfully!\n\nID : %v\nCreated At : %v\nUpdated At : %v\nName of feed : %v\nURL : %v\nUserID : %v\n", feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
	return nil
}
