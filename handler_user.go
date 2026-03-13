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

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("Succesfully cleared users table!\n")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, user := range users {
		if s.cfg.Current_user_name == user.Name {
			fmt.Printf("- %s (current)\n", user.Name)
		} else {
			fmt.Printf("- %s\n", user.Name)
		}

	}

	return nil
}

func handlerAggregate(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"

	rssFeed, err := fetchFeed(context.Background(), url)

	if err != nil {
		return err
	}

	fmt.Printf("%v", *rssFeed)
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Insufficient args!\n")
	}

	current_user := s.cfg.Current_user_name
	url := cmd.args[0]

	userID := user.ID

	feed, err := s.db.GetFeedFromURL(context.Background(), url)

	if err != nil {
		return err
	}

	feedFollowItem := database.CreatedFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedID:    feed.ID,
		UserID:    userID,
	}

	_, err = s.db.CreatedFeedFollow(context.Background(), feedFollowItem)

	if err != nil {
		return err
	}

	fmt.Printf("Name of Feed : %s\nUser : %s\n", feed.Name, current_user)
	return nil

}

func handlerFollowing(s *state, cmd command, user database.User) error {
	current_user := s.cfg.Current_user_name

	feed_follow, err := s.db.GetFeedFollowsForUser(context.Background(), current_user)

	if err != nil {
		return err
	}

	fmt.Printf("Feeds followed by %s:\n", current_user)
	for _, row := range feed_follow {
		fmt.Printf("- %s\n", row.FeedName)
	}

	return nil
}
