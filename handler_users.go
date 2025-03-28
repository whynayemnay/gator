package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/whynayemnay/gator/internal/database"
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

func handlerRegister(state *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("register command needs one argument: username")
	}

	name := cmd.arguments[0]

	user, err := state.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})

	if err != nil {
		// Check for duplicate entry error (unique constraint violation)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			fmt.Println("Error: user already exists")
			os.Exit(1) // Exit with code 1
		}
		// Handle other potential errors
		log.Fatalf("Failed to create user: %v", err)
	}

	state.state.SetUser(user.Name)
	fmt.Println("User created:", user.Name)

	return nil
}

func handlerLogins(state *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("command was empty")
	}

	username := cmd.arguments[0]
	err := state.state.SetUser(username)
	if err != nil {
		return err
	}

	user, err := state.db.GetUser(context.Background(), username)
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
