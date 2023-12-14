package storage

import (
	"HHBot/models"
	"context"
)

type Storage interface {
	GetVacancies(context.Context, models.Filter, string) ([]models.Vacancy, error)
	SetSettings(context.Context, models.Filter) error
	GetSettings(context.Context, int) (*models.Filter, error)
}
