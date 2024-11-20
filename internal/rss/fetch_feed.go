package rss

import (
	"context"
	"encoding/xml"
	"html"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var feeds RSSFeed
	decoder := xml.NewDecoder(res.Body)
	if err = decoder.Decode(&feeds); err != nil {
		return nil, err
	}
	feeds = cleanFeed(feeds)

	return &feeds, nil
}

func cleanFeed(feed RSSFeed) RSSFeed {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i] = item
	}
	return feed
}
