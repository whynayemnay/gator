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
