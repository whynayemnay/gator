package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func handlerLogins(state *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("command was empty")
	}

	username := cmd.arguments[0]
	err := state.state.SetUser(username)
	if err != nil {
		return err
	}

	user, err := state.db.GetUsers(context.Background(), username)
	if err != nil {
		fmt.Println("erro user does not exist")
		os.Exit(1)
	}

	err = state.state.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Println("logged with the name", user.Name)
	return nil
}
