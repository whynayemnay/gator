package main

import (
	"context"
	"fmt"
)

func handlerReset(state *state, cmd command) error {

	err := state.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users %w", err)
	}
	fmt.Println("Deleted all users in the user table!")
	return nil
}
