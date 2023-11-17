package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	AppPort    string
	Host       string
	TokenHH    string
	TokenTG    string
}

func GetConfig(path string) (*Config, error) {
	config := &Config{}
	if err := godotenv.Load(path); err != nil {
		return config, err
	}

	config.DBUser = os.Getenv("DB_USER")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBName = os.Getenv("DB_NAME")
	config.DBPort = os.Getenv("DB_PORT")
	config.AppPort = os.Getenv("APP_PORT")
	config.Host = os.Getenv("HOST")
	config.TokenHH = os.Getenv("TOKEN_HH")
	config.TokenTG = os.Getenv("TOKEN_TG")

	return config, nil
}
