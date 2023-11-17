package main

import (
	"log"

	"github.com/MithraRa/hh-bot/config"
	repository "github.com/MithraRa/hh-bot/internal/repository/postgres"
)

func main() {
	cfg, err := config.GetConfig("../.env")
	if err != nil {
		log.Fatal(err)
	}

	_, err = repository.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

}
