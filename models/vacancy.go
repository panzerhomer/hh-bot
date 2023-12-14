package models

type Filter struct {
	UserID     int
	Salary     string
	Experience string
	City       string
}

type Vacancy struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Salary     string `json:"salary"`
	City       string `json:"city"`
	Company    string `json:"company"`
	Experience string `json:"experience"`
	Skills     string `json:"skills"`
	URL        string `json:"url"`
}
