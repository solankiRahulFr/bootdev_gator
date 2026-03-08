package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"gator/internal/database"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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

// Function to fetch the RSS feed from the url
func fetchFeed(url string) (*RSSFeed, error) {
	ctx := context.Background()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	// setting the user-agent to gator to avoid 403 forbidden error from some sites
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, err
	}
	// Unescape channel fields
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Unescape each item
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

// function to get the next feed to fetch from the database and scrape it
func scrapeFeed(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Printf("no feed to fetch: %v\n", err)
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		fmt.Printf("mark feed fetched: %v\n", err)
		return err
	}

	rssFeed, err := fetchFeed(feed.Url)
	if err != nil {
		return err
	}
	// Print the entire struct
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("Creating posts in the Post DB")
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: time.Now().UTC(),
			FeedID:      feed.ID,
		})
		if err != nil {
			fmt.Printf("create post: %v\n", err)
			continue
		}
	}
	return nil
}
