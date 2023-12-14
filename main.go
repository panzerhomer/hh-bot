package main

import (
	tgClient "HHBot/clients/telegram"
	"HHBot/config"
	event_consumer "HHBot/consumer/event-consumer"
	"HHBot/events/telegram"
	"HHBot/storage/postgres"
	"fmt"
	"log"
)

const (
	batchSize = 100
)

func main() {
	cfg, err := config.GetConfig(".env")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.DBName)

	s, err := postgres.New(dsn)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(cfg.Bot.Host, cfg.Bot.Token),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
