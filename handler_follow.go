package main

import (
	"context"
	"fmt"
	"time"

	"github.com/SoFarCalm/blobgator/internal/database"
	"github.com/google/uuid"
)

func handlerUnfollow(s *state, cmd command, user database.GetUserRow) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Please provide a feed url")
	}

	feed, err := s.db.GetFeedName(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	s.db.DeleteFeedFollows(context.Background(), database.DeleteFeedFollowsParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("User %s has unfollowed %s \n", user.Name, feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.GetUserRow) error {

	following, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		fmt.Errorf("Error fetching follows: %w", err)
	}

	for _, follow := range following {
		fmt.Printf("- %s\n", follow.FeedName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.GetUserRow) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("Please enter a URL")
	}

	url := cmd.Args[0]

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
	fmt.Printf(" * Feed:    %v\n", feed_follow.FeedName)
	fmt.Printf(" * User:    %v\n", feed_follow.UserName)
}
