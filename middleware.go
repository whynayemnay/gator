package main

import (
	"context"
	"fmt"

	"github.com/whynayemnay/gator/internal/database"
)

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Retrieve the user based on the session or context
		user, err := s.db.GetUser(context.Background(), s.state.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to fetch the user info: %w", err)
		}

		// Call the original handler with the user parameter
		return handler(s, cmd, user)
	}
}
