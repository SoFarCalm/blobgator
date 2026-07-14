package main

import (
	"context"
	"fmt"
	"strconv"
)

func handlerBrowse(s *state, cmd command) error {
	limit := int32(2)

	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			fmt.Errorf("Please set a valid integer limit %w", err)
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPosts(context.Background(), limit)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("-Title:       %s\n", post.Title.String)
		fmt.Printf("-Description: %s\n", post.Description.String)
		fmt.Printf("-Url:         %s\n", post.Url.String)
		fmt.Printf("-Description: %v\n", post.PublishedAt.Time)
		fmt.Println()
	}

	return nil

}
