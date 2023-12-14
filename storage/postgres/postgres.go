package postgres

import (
	"HHBot/models"
	"HHBot/utils"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetVacancies(ctx context.Context, f models.Filter, search string) ([]models.Vacancy, error) {
	querySelect := createQuerySelect(f, search)

	rows, err := s.db.QueryContext(ctx, querySelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []models.Vacancy
	for rows.Next() {
		var v models.Vacancy
		err = rows.Scan(&v.ID, &v.Title, &v.Salary, &v.City, &v.Company, &v.Experience, &v.Skills, &v.URL)
		if err != nil {
			return nil, err
		}
		vacancies = append(vacancies, v)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	log.Print("[storage]", vacancies, querySelect)

	return vacancies, nil
}

func (s *Storage) SetSettings(ctx context.Context, f models.Filter) error {
	const errorMsg = "can't set settings"
	queryInsert := "INSERT INTO users VALUES ($1, $2, $3, $4)"

	if _, err := s.db.ExecContext(ctx, queryInsert, f.UserID, f.City, f.Salary, f.Experience); err != nil {
		return utils.Wrap(errorMsg, err)
	}

	return nil
}

func (s *Storage) GetSettings(ctx context.Context, UserID int) (*models.Filter, error) {
	querySelect := "SELECT id, experience, salary, city FROM users"

	rows, err := s.db.QueryContext(ctx, querySelect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var f *models.Filter

	err = rows.Scan(&f.UserID, &f.Experience, &f.Salary, &f.City)
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return f, nil
}

func createQuerySelect(f models.Filter, search string) string {
	condition := "'%" + search + "%'"
	query := "SELECT id, title, salary, city, company, experience, skills, url FROM vacancies WHERE title LIKE " + condition
	query = "SELECT id, title, salary, city, company, experience, skills, url FROM vacancies"
	// if f.City != "" {
	// 	query += " AND city = " + "Москва"
	// }
	// if f.Experience != "" {
	// 	query += " AND experience = " + f.Experience
	// }
	// if f.Salary != "" {
	// 	query += " AND salary >= " + f.Salary
	// }

	return query + " limit 1"
}

// func (s *Storage) Init() error {

// }
