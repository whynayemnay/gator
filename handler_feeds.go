package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/whynayemnay/gator/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error during the request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error during the response: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error during the data reading: %w", err)
	}
	var rss RSSFeed
	err = xml.Unmarshal(data, &rss)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error during the unmarshalig of data: %w", err)
	}
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(rss.Channel.Item[i].Title)
		rss.Channel.Item[i].Description = html.UnescapeString(rss.Channel.Item[i].Description)
	}

	return &rss, nil
}

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	feedURL := "https://www.wagslane.dev/index.xml"

	rss, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Printf("Fetched RSS Feed:\n%+v\n", rss)
	return nil
}

func handlerFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 2 {
		return fmt.Errorf("give two arguments 1: feed name, 2: feed url")
	}
	feedName := cmd.arguments[0]
	feedURL := cmd.arguments[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating a new feed entry: %w", err)
	}
	fmt.Printf("new feed added to db: \n"+
		"feedID: %v\n"+
		"created at: %v\n"+
		"updated at: %v\n"+
		"name: %v\n"+
		"url: %v\n"+
		"user_id: %v\n", feed.ID, feed.CreatedAt, feed.UpdatedAt,
		feed.Name, feed.Url, feed.UserID)

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("erorr inserting into the feed_follow table: %w", err)
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("command doesn't need any arguments")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}
	for i, feed := range feeds {
		fmt.Printf("Feed #%v\n"+
			"feed name: %v\n"+
			"feed url: %v\n"+
			"user who added feed: %v\n", i+1, feed.Name, feed.Url, feed.Name_2)
	}

	return nil
}

func handlerInsertFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("provide feed name as argument")
	}
	url := cmd.arguments[0]
	feed, err := s.db.SelectFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("error fetching the feed data by url: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("erorr inserting into the feed_follow table: %w", err)
	}

	fmt.Printf("Created a new follow entry for:\n"+
		"user: %v\n"+
		"feed name: %v\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 0 {
		return fmt.Errorf("error to manny arguments passed for the command")
	}
	query, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error getting follows for the user")
	}
	fmt.Println("user is currently following:")
	for _, feed := range query {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handlerDeleteFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("provide 1 argument: url of the feed")
	}
	feed, err := s.db.SelectFeedByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't delte the feed entry: %w", err)
	}
	return nil
}
