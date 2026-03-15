package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/roshankaranth/gator/internal/database"
)

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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Insufficient args!")
	}

	user_id := user.ID
	feed, err := s.db.GetFeedFromURL(context.Background(), cmd.args[0])

	if err != nil {
		return err
	}

	feed_id := feed.ID

	feed_follow := database.DeleteFeedFollowParams{
		UserID: user_id,
		FeedID: feed_id,
	}

	err = s.db.DeleteFeedFollow(context.Background(), feed_follow)

	if err != nil {
		return err
	}

	fmt.Printf("Unfollowed %s successfully!\n", feed.Name)
	return nil
}
