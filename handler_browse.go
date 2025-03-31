package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/whynayemnay/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.arguments) == 1 {
		if specifiedLimit, err := strconv.Atoi(cmd.arguments[0]); err == nil {
			limit = specifiedLimit
		} else {
			return fmt.Errorf("invalid limit: %w", err)
		}
	}

	posts, err := s.db.GetPostUser(context.Background(), database.GetPostUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("couldn't get posts for user: %w", err)
	}

	fmt.Printf("found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Local().Format("2006 Jan 2"), post.Name)
		fmt.Printf("--- %s ---\n", post.Title)

		desc := "(No description available)!!!!"
		if post.Description.Valid {
			desc = post.Description.String
		}
		fmt.Printf("    %s\n", desc)

		// Ensure URL is printed correctly
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
