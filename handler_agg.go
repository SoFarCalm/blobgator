package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/SoFarCalm/blobgator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("Please provide a latency for requests e.g('1s', '1m', '1h')")
	}

	timeToSleep, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeToSleep)

	ticker := time.NewTicker(timeToSleep)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	notWorking := []string{"BootDev", "reddit-go", "julia-evans", ""}

	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Ran GetNextFeedToFetch...")

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	isNotWorking := slices.Contains(notWorking, feed.Name)
	if isNotWorking == true {
		fmt.Printf("Feed %s is not working, skipped\n", feed.Name)
		return nil
	}

	fmt.Printf("Fetching feed %s\n", feed.Name)

	fetchedFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	for _, post := range fetchedFeed.Channel.Item {
		fmt.Println("entering for loop here")
		parsedPubTime, err := time.Parse(time.RFC3339, post.PubDate)
		fmt.Println(parsedPubTime)
		if err != nil {
			return err
		}
		fmt.Println("Creating post...")
		postCreated, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: sql.NullTime{
				Time:  time.Now().UTC(),
				Valid: !time.Now().UTC().IsZero(),
			},
			Title: sql.NullString{
				String: post.Title,
				Valid:  post.Title != "",
			},
			Url: sql.NullString{
				String: post.Link,
				Valid:  post.Link != "",
			},
			Description: sql.NullString{
				String: post.Description,
				Valid:  post.Description != "",
			},
			PublishedAt: sql.NullTime{
				Time:  parsedPubTime,
				Valid: !parsedPubTime.IsZero(),
			},
			FeedID: feed.ID,
		})
		if err != nil {
			fmt.Errorf("Could not create post %w", err)
		}
		fmt.Println(postCreated)
	}

	// for _, v := range fetchedFeed.Channel.Item {
	// 	fmt.Printf(" -%s\n", v)
	// }

	return nil
}
