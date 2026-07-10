package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SoFarCalm/blobgator/internal/database"
	"github.com/google/uuid"
)

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Errorf("Error fecthing feeds: %w", err)
	}

	for _, feed := range feeds {
		createdBy, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			fmt.Errorf("Error fetching username")
		}
		printFeed(feed)
		fmt.Printf(" * User:      %v\n", createdBy)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %v <name>, example.com <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]
	currentUser := s.cfg.CurrentUserName

	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		return fmt.Errorf("Error fecthing users")
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})

	printFeed(feed)

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:      %v\n", feed.ID)
	fmt.Printf(" * Name:    %v\n", feed.Name)
	fmt.Printf(" * URL:     %v\n", feed.Url)
}
