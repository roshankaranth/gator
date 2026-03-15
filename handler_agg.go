package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/roshankaranth/gator/internal/database"
)

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("Insufficient args!")
	}

	time_between_reqs, err := time.ParseDuration(cmd.args[0])

	if err != nil {
		return err
	}

	fmt.Printf("Collecting feed every %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

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
