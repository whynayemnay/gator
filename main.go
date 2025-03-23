package main

import (
	"fmt"
	"log"
	"os"

	"github.com/whynayemnay/gator/internal/config"
)

type state struct {
	state *config.Config
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config:", err)
	}
	fmt.Println("DBURL config file: ", cfg)

	appState := &state{state: &cfg}

	cmds := commands{command: make(map[string]func(*state, command) error)}

	cmds.register("login", handlerLogins)

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
