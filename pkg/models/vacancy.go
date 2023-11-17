package models

import "strconv"

type Vacancies struct {
	Items []Vacancy `json:"items"`
}

type Vacancy struct {
	Title    string   `json:"title"`
	Salary   Salary   `json:"salary"`
	Area     Area     `json:"area"`
	Employer Employer `json:"employer"`
	URL      string   `json:"url"`
}

type Area struct {
	City string `json:"name"`
}

type Salary struct {
	From     int    `json:"from"`
	To       int    `json:"to"`
	Currency string `json:"currency"`
}

type Employer struct {
	Name string `json:"name"`
}

func ParseSalary(salary Salary) string {
	from := strconv.Itoa(salary.From)
	if from == "0" {
		from = ""
	} else {
		from = "от " + from + " "
	}

	to := strconv.Itoa(salary.To)
	if to == "0" {
		to = ""
	} else {
		to = "до " + to + " "
	}

	if from == "" && to == "" {
		return "Зарплата: Не указана"
	}

	return "Зарплата: " + from + to + salary.Currency
}

func ParseExperience(req string) string {
	switch req {
	case "Нет опыта":
		return "noExperience"
	case "От 1 до 3 лет":
		return "between1And3"
	case "От 3 до 6 лет":
		return "between3And6"
	case "Более 6 лет":
		return "moreThan6"
	default:
		return ""
	}
}
