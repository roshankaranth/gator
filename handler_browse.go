package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/roshankaranth/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	lim := int32(2)
	if len(cmd.args) >= 1 {
		v, err := strconv.Atoi(cmd.args[0])

		if err != nil {
			return err
		}

		lim = int32(v)

	}

	GetPostArgs := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  lim,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), GetPostArgs)

	if err != nil {
		return err
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt, post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil

}
