package main

import (
	"log"

	"botjob.com/config"
	repository "botjob.com/internal/repository/postgres"
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
