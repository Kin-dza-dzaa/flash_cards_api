package main

import (
	"log"

	"github.com/Kin-dza-dzaa/flash_cards_api/config"
	"github.com/Kin-dza-dzaa/flash_cards_api/internal/app"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
