package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/whynayemnay/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) < 1 || len(cmd.arguments) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Println("problem parsing the date")
		}
		err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "unique constainr vialation, URL was already added") {
				continue
			}
			log.Printf("couldn'T add post to DB %v", err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
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
