package repository

import (
	"database/sql"
	"fmt"

	"botjob.com/config"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(config *config.Config) (*Repository, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.DBName,
		config.DBPassword,
		config.Host,
		config.DBPort,
		config.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}
