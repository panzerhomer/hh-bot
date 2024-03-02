package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		Bot      Bot
		Postgres Postgres
	}

	Bot struct {
		Name  string
		Token string
		Host  string
	}

	Postgres struct {
		Host     string
		Port     string
		Username string
		Password string
		DBName   string
		SSL      string
	}
)

func GetConfig(path string) (*Config, error) {
	config := &Config{}
	if err := godotenv.Load(path); err != nil {
		return config, err
	}

	bot := Bot{
		Name:  os.Getenv("BOT_NAME"),
		Token: os.Getenv("BOT_TGTOKEN"),
		Host:  os.Getenv("BOT_HOST"),
	}
	config.Bot = bot

	posgres := Postgres{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USERNAME"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DBNAME"),
		SSL:      os.Getenv("POSTGRES_SSL_MODE"),
	}
	config.Postgres = posgres

	return config, nil
}
