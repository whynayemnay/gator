package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/whynayemnay/gator/internal/config"
	"github.com/whynayemnay/gator/internal/database"
)

type state struct {
	db    *database.Queries
	state *config.Config
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config:", err)
	}
	// fmt.Println("DBURL config file: ", cfg)

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("error opening the postgres connection")
	}

	dbQueries := database.New(db)

	appState := &state{state: &cfg}

	appState.db = dbQueries

	cmds := commands{command: make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogins)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerInsertFeed)
	cmds.register("following", handlerFollowing)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("error: no command.")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	cmd := command{name: cmdName, arguments: cmdArgs}

	if err := cmds.run(appState, cmd); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
