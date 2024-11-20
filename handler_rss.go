package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/database"
	"gator/internal/rss"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(state *State, cmd Command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("incorrect command usage: %s <interval>", cmd.name)
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("cannot parse time: %v", err)
	}
	ticker := time.NewTicker(duration)
	for range ticker.C {
		err = scrapeFeeds(state.db)
		if err != nil {
			fmt.Printf("Error trying to scrape feed: %v\n", err)
		}
	}
	return nil
}

func handlerAddFeed(state *State, cmd Command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("expect 2 arguemtn, found %d", len(cmd.args))
	}

	feed, err := state.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %v", err)
	}
	fmt.Println("Successfully added feed")
	return nil
}

func handlerFeeds(state *State, _ Command) error {
	feeds, err := state.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("fail to get feed: %v", err)
	}
	for _, feed := range feeds {
		fmt.Printf("\t - id: %v\n", feed.ID)
		fmt.Printf("\t   created at: %v\n", feed.CreatedAt)
		fmt.Printf("\t   updated at: %v\n", feed.UpdatedAt)
		fmt.Printf("\t   name: %v\n", feed.Name)
		fmt.Printf("\t   url: %v\n", feed.Url)
		fmt.Printf("\t   username: %v\n", feed.Username)
	}
	return nil
}

func handlerFollow(state *State, cmd Command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("incorrect command usage: %s <url>", cmd.name)
	}
	feed, err := state.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("feed not found: %v", err)
	}

	follow, err := state.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("cannot follow feed: %v", err)
	}
	fmt.Printf("User %v followed %v succesffuly!\n", follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing(state *State, _ Command, user database.User) error {
	feeds, err := state.db.GetFeedFollowsByUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("cannot get feeds: %v", err)
	}

	fmt.Printf("Here are all the feeds followed by %v:\n", user.Name)
	if len(feeds) == 0 {
		fmt.Println("\tNone")
		return nil
	}
	for _, feed := range feeds {
		fmt.Printf("\t - id: %v\n", feed.ID)
		fmt.Printf("\t   created at: %v\n", feed.CreatedAt)
		fmt.Printf("\t   updated at: %v\n", feed.UpdatedAt)
		fmt.Printf("\t   name: %v\n", feed.Name)
		fmt.Printf("\t   url: %v\n", feed.Url)
		fmt.Printf("\t   created by: %v\n", feed.CreatedBy)
	}

	return nil
}

func handlerUnfollow(state *State, cmd Command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("incorrect command usage: %s <url>", cmd.name)
	}
	feed, err := state.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("cannot find feed: %v", err)
	}
	err = state.db.DeleteFollowFeed(context.Background(), database.DeleteFollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow: %v", err)
	}

	fmt.Println("Unfollowed successfully")
	return nil
}

func handlerBrowse(state *State, cmd Command, user database.User) error {
	var limit int32 = 2
	if len(cmd.args) > 0 {
		uLimit, err := strconv.ParseInt(cmd.args[0], 0, 0)
		if err != nil {
			return fmt.Errorf("unable to parse int: %v", err)
		}
		if uLimit < 1 {
			return fmt.Errorf("limit must be >= 1")
		}
		limit = int32(uLimit)
	}

	posts, err := state.db.GetPostForUser(context.Background(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("cannot get posts: %v", err)
	}

	// fmt.Printf("Title: %v\nDescription: %v\nLink: %v\n", data.Channel.Title, strings.TrimSpace(data.Channel.Description), data.Channel.Link)
	for _, post := range posts {

		fmt.Printf("\n\t - Title: %v\n\t   Description: %v\n\t   Publication date: %v\n\t   Link: %v\n", post.Title, strings.TrimSpace(post.Description), post.PublishedAt, post.Url)
	}

	return nil
}

func scrapeFeeds(db *database.Queries) error {
	feed, err := db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving next feed to scrape: %v", err)
	}

	err = db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     time.Now(),
		ID:            feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error updating feed as fetched: %v", err)
	}

	data, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed: %v", err)
	}

	postNo := 0
	dupNo := 0
	for _, item := range data.Channel.Item {
		layout := "Mon, 02 Jan 2006 15:04:05 -0700"

		pubDate, err := time.Parse(layout, item.PubDate)
		if err != nil {
			fmt.Printf("Cannot parse pubDate: %v\n", err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			FeedID:      feed.ID,
			PublishedAt: pubDate,
		})
		if err != nil {
			var pqErr *pq.Error
			if ok := errors.As(err, &pqErr); ok && pqErr.Code == "23505" {
				dupNo++
			} else {
				fmt.Printf("Err while saving post: %v\n", err)
			}
		} else {
			postNo++
		}
	}
	fmt.Printf("Scrape %v, added %v posts, duplicate %v posts\n", feed.Name, postNo, dupNo)
	// fmt.Printf("Title: %v\nDescription: %v\nLink: %v\n", data.Channel.Title, strings.TrimSpace(data.Channel.Description), data.Channel.Link)
	// for _, item := range data.Channel.Item {
	// 	fmt.Printf("\n\t - Title: %v\n\t   Description: %v\n\t   Publication date: %v\n\t   Link: %v\n", item.Title, strings.TrimSpace(item.Description), item.PubDate, item.Link)
	// }

	return nil
}
