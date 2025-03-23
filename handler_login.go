package main

import (
	"errors"
	"fmt"
)

func handlerLogins(state *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("command was empty")
	}

	err := state.state.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("logged with the name", cmd.arguments[0])
	return nil
}
