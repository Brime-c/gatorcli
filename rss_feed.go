package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

// RSSFeed models the top-level structure of an RSS 2.0 XML feed.
// It maps the XML root elements to Go structs.
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// RSSItem models an individual article, post, or entry within an RSS feed channel.
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// fetchFeed requests an RSS feed from the given URL, parses the XML response,
// and decodes it into an RSSFeed struct. It unescapes HTML entities in titles
// and descriptions before returning.
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// 1. Build the HTTP request with the provided context for cancellation/timeouts
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}

	// 2. Set a polite User-Agent so feed hosts know who is fetching their data
	req.Header.Set("User-Agent", "gator")

	// 3. Execute the network request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 4. Read the entire raw XML body from the response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 5. Unmarshal (parse) the XML data into our Go structs
	var feed RSSFeed
	err = xml.Unmarshal(bodyBytes, &feed)
	if err != nil {
		return nil, err
	}

	// 6. Clean up HTML entities (e.g., converting "&" back to "&")
	// to make the content readable in the terminal
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}
