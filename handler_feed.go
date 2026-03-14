package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/roshankaranth/gator/internal/database"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.args) < 2 {
		return fmt.Errorf("insufficient arguments!\n")
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

	feedFollowItem := database.CreatedFeedFollowParams{
		ID:        user_feed.ID,
		CreatedAt: user_feed.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		FeedID:    feed.ID,
		UserID:    userID,
	}

	_, err = s.db.CreatedFeedFollow(context.Background(), feedFollowItem)

	if err != nil {
		return err
	}

	fmt.Printf("Feed created succesfully!\n\nID : %v\nCreated At : %v\nUpdated At : %v\nName of feed : %v\nURL : %v\nUserID : %v\n", feed.ID, feed.CreatedAt, feed.UpdatedAt, feed.Name, feed.Url, feed.UserID)
	return nil
}

func handlerFeeds(s *state, cmd command) error {

	feeds, err := s.db.GetFeeds(context.Background())

	if err != nil {
		return err
	}

	fmt.Printf("Feeds : \n\n")
	for i, feed := range feeds {
		user, err := s.db.GetUserFromID(context.Background(), feed.UserID)

		if err != nil {
			return err
		}
		fmt.Printf("%v)", i+1)
		fmt.Printf(" Name : %s\n", feed.Name)
		fmt.Printf("   URL : %s\n", feed.Url)
		fmt.Printf("   User : %s\n\n", user.Name)

	}

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

func scrapeFeeds(s *state) error {
	feed_to_fetch, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return err
	}

	mark_feed_params := database.MarkFeedFetchedParams{
		ID: feed_to_fetch.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		},
	}

	err = s.db.MarkFeedFetched(context.Background(), mark_feed_params)

	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), feed_to_fetch.Url)

	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {

		pubDate, err := parseTime(item.PubDate)

		if err != nil {
			return err
		}

		post := database.CreatePostsParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: pubDate,
			FeedID:      feed_to_fetch.ID,
		}
		_, err = s.db.CreatePosts(context.Background(), post)

		if err != nil {
			return err
		}
	}
	return nil

}

func parseTime(s string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse time")
}
