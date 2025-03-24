package main

import (
	"context"
	"fmt"
	"os"
)

func handlerReset(state *state, cmd command) error {

	err := state.db.DeleteUsers(context.Background())
	if err != nil {
		fmt.Println("error deleting user: ", err)
		os.Exit(1)
	}
	fmt.Println("Deleted all users in the user table!")
	return nil
}
