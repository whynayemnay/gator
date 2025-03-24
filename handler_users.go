package main

import (
	"context"
	"fmt"
)

func handlerGetUsers(state *state, cmd command) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting all users %w", err)
	}
	for _, user := range users {
		if user == state.state.CurrentUserName {
			fmt.Printf("* %s (current)\n", user)
			continue
		}
		fmt.Println("*", user)
	}
	return nil
}
