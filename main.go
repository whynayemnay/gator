package main

import (
	"fmt"
	"log"

	"github.com/whynayemnay/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("error reading config:", err)
	}
	fmt.Println("DBURL config file: ", cfg)

	err = cfg.SetUser("whynay")
	if err != nil {
		log.Fatal("error setting user:", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatal("errro reading the file: ", err)
	}

	fmt.Println("DB URL:", cfg.DBURL)
	fmt.Println("Current user name: ", cfg.CurrentUserName)
}
