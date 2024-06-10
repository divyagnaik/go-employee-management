package handler

import (
	"personal/go-employee-management/models"
	"personal/go-employee-management/pagination"

	"gofr.dev/pkg/gofr"
)

type employeeService interface {
	Create(ctx *gofr.Context, e *models.Employee) (*models.Employee, error)
	Get(ctx *gofr.Context, id string) (*models.Employee, error)
	GetAll(ctx *gofr.Context, page *pagination.Page) ([]models.Employee, error)
	Update(ctx *gofr.Context, e *models.Employee) (*models.Employee, error)
	Delete(ctx *gofr.Context, id string) error
}
