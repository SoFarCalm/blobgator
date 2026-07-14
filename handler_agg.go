package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
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
		//parsedPubTime, err := time.Parse(time.RFC3339, post.PubDate)
		pubDate, err := parseFeedDate(post.PubDate)
		if err != nil {
			return err
		}

		fmt.Println(pubDate)
		if err != nil {
			return err
		}
		fmt.Println("Creating post...")
		now := time.Now().UTC()
		postCreated, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: sql.NullTime{
				Time:  now,
				Valid: true,
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
				Valid:  post.Link != "",
			},
			PublishedAt: sql.NullTime{
				Time:  pubDate,
				Valid: !pubDate.IsZero(),
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

func parseFeedDate(dateStr string) (time.Time, error) {
	// Clean whitespace around the incoming string
	dateStr = strings.TrimSpace(dateStr)

	// List of common formats used across various RSS, Atom, and custom JSON feeds
	formats := []string{
		time.RFC3339,          // Atom / JSON Feeds (e.g. 2026-07-13T21:22:00Z)
		time.RFC1123Z,         // RSS 2.0 with numeric zones (e.g. Mon, 13 Jul 2026 21:22:00 +0000)
		time.RFC1123,          // RSS 2.0 with named zones (e.g. Mon, 13 Jul 2026 21:22:00 GMT)
		"2006-01-02 15:04:05", // Legacy Postgres style fallback timestamps
	}

	for _, format := range formats {
		if parsedTime, err := time.Parse(format, dateStr); err == nil {
			// Convert to UTC explicitly so it stores correctly in your PostgreSQL database
			return parsedTime.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date string: %s", dateStr)
}
