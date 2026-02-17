package main

import (
	"log"

	"github.com/josinaldojr/imobifix-api/internal/app"
	"github.com/josinaldojr/imobifix-api/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
