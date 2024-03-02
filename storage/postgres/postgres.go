package postgres

import (
	"HHBot/models"
	"HHBot/utils"
	"context"
	"database/sql"
	"fmt"

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

func (s *Storage) GetVacancies(ctx context.Context, f *models.Filter, search string) ([]models.Vacancy, error) {
	query := createQuerySelect(f, search)

	rows, err := s.db.QueryContext(ctx, query)
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
		if v.Salary == "nan" {
			v.Salary = "з/п не указана"
		} else if v.Salary != "з/п не указана" {
			v.Salary += " ₽"
		}
		vacancies = append(vacancies, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return vacancies, nil
}

func (s *Storage) SetSettings(ctx context.Context, f *models.Filter) error {
	const errorMsg = "can't set settings"

	query := ""
	filter, _ := s.GetSettings(ctx, f.UserID)
	if filter.UserID == 0 {
		query = "INSERT INTO users VALUES ($1, $2, $3, $4)"
	} else {
		query = "UPDATE users SET city = $1, salary = $2, experience = $3 WHERE id = $4"
	}

	if _, err := s.db.ExecContext(ctx, query, f.UserID, f.City, f.Salary, f.Experience); err != nil {
		return utils.Wrap(errorMsg, err)
	}

	return nil
}

func (s *Storage) GetSettings(ctx context.Context, UserID int) (models.Filter, error) {
	query := "SELECT id, city, salary, experience FROM users where id = $1"

	var f models.Filter

	rows := s.db.QueryRowContext(ctx, query, UserID)

	err := rows.Scan(&f.UserID, &f.City, &f.Salary, &f.Experience)
	if err != nil {
		return f, err
	}
	if err = rows.Err(); err != nil {
		return f, err
	}

	return f, nil
}

func createQuerySelect(f *models.Filter, search string) string {
	condition := "'%" + search + "%'"
	query := "SELECT id, title, salary, city, company, experience, skills, url FROM vacancies WHERE title LIKE  " + condition

	if f.City != "" {
		query += fmt.Sprintf(" AND city = '%s' ", f.City)
	}
	if f.Experience != "" {
		query += fmt.Sprintf(" AND experience = '%s' ", f.Experience)
	}
	// if f.Salary != "nan" && f.Salary != "з/п не указана" {
	// 	salary, err := strconv.Atoi(f.Salary)
	// 	if err != nil {

	// 	} else {
	// 		query += " AND salary >= " + salary
	// 	}
	// }

	return query + " limit 3"
}
