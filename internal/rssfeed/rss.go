package rssfeed

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// What context.Context means in client.NewRequestWithContext()
// In Go, context.Context is a standard way to carry request-scoped metadata, deadlines, and cancellation signals through a call chain.

// When you call something like:

// the ctx parameter means:

// ctx is the operation context for that HTTP request
// it can carry a deadline or timeout
// it can be canceled
// it is propagated to lower-level operations such as network I/O
// Why it matters
// If ctx is canceled or its deadline expires:

// the request can be aborted
// resources can be released
// the client library can stop waiting for a response
// Example use cases:

// request timeout in a web handler
// canceling an API call when the caller disconnects
// enforcing a global deadline across multiple requests
// Typical pattern
// This means the created request will automatically stop if:

// 5 seconds pass
// cancel() is called
// the parent's context is canceled
// Key points
// context.Context is not for passing optional arguments like config values
// it should be the first parameter in function signatures when present
// client.NewRequestWithContext() uses it to attach cancellation/deadline behavior to the request

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (c *Client) FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		fmt.Println("Error creating request")
		return &RSSFeed{}, nil
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(("Error fecthing RSS feed"))
		return &RSSFeed{}, nil
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading xml data")
		return &RSSFeed{}, nil
	}

	var feed RSSFeed

	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return &RSSFeed{}, nil
	}

	fmt.Println("Maybe this worked...")
	fmt.Println(feed)
	return &feed, nil
}
