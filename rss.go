package main

import (
	"fmt"
	"context"
	"encoding/xml"
	"net/http"
	"time"
	"io"
	"html"
	"log"
	"database/sql"

	"github.com/google/uuid"
	"github.com/gclinoz/Gator-go/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title		string		`xml:"title"`
		Link		string		`xml:"link"`
		Description string		`xml:"description"`
		Item		[]RSSItem	`xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title		string	`xml:"title"`
	Link		string	`xml:"link"`
	Description	string	`xml:"description"`
	PubDate		string	`xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	data := RSSFeed{}
	err = xml.Unmarshal(body, &data)
	if err != nil {
		return &RSSFeed{}, err
	}

	data.Channel.Title = html.UnescapeString(data.Channel.Title)
	data.Channel.Description = html.UnescapeString(data.Channel.Description)
	for i, val := range data.Channel.Item {
		val.Title = html.UnescapeString(val.Title)
		val.Description = html.UnescapeString(val.Description)
		data.Channel.Item[i] = val
	}

	return &data, nil
}

func scrapeFeeds(s *state) {
	ctx := context.Background()

	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Printf("couldn't get next feed to fetch: %w", err)
	}

	_, err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s fetched: %v", nextFeed.Name, err)
	}

	data, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		log.Printf("couldn't collect feed %s: %v", nextFeed.Name, err)
	}

	for _, val := range data.Channel.Item {
		parsed, err := time.Parse(time.RFC1123Z, val.PubDate)
		if err != nil {
			parsed = time.Time{}
		}

		postParam := database.CreatePostParams{
			ID:				uuid.New(),
			CreatedAt:		time.Now(),
			UpdatedAt:		time.Now(),
			Title:			sql.NullString{String: val.Title, Valid: true},
			Url:			val.Link,
			Description:	sql.NullString{String: val.Description, Valid: true},
			PublishedAt:	sql.NullTime{Time: parsed, Valid: true},
			FeedID:			nextFeed.ID,
		}
		_, err = s.db.CreatePost(ctx, postParam)
		if err != nil {
			log.Printf("couldn't insert posts from %s: %v", nextFeed.Name, err)
		}
	}
	fmt.Printf("Feed %s collected, %v posts stored\n", nextFeed.Name, len(data.Channel.Item))
}
