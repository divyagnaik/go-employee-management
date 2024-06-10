package service

import (
	"personal/go-employee-management/models"
	"personal/go-employee-management/pagination"

	"gofr.dev/pkg/gofr"
)

type employeeStore interface {
	Create(ctx *gofr.Context, e *models.Employee) (*models.Employee, error)
	Get(ctx *gofr.Context, id int64) (*models.Employee, error)
	GetAll(ctx *gofr.Context, page *pagination.Page) ([]models.Employee, error)
}
