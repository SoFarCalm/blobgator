package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SoFarCalm/blobgator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Please enter a URL")
	}

	url := cmd.Args[0]
	currentUser := s.cfg.CurrentUserName

	user, err := s.db.GetUser(context.Background(), currentUser)
	if err != nil {
		fmt.Errorf("Error fecthing user name: %w", err)
	}

	feed, err := s.db.GetFeedName(context.Background(), url)
	if err != nil {
		fmt.Errorf("Error fetching feed name: %w", err)
	}

	feed_follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	printFeedFollow(feed_follow)

	return nil
}

func printFeedFollow(feed_follow database.CreateFeedFollowRow) {
	fmt.Printf(" * Feed:      %v\n", feed_follow.FeedName)
	fmt.Printf(" * User:    %v\n", feed_follow.UserName)
}
