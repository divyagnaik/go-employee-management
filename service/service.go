package service

import (
	"fmt"
	"personal/go-employee-management/models"
	"personal/go-employee-management/pagination"
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

type Service struct {
	store employeeStore
}

func New(e employeeStore) *Service {
	return &Service{store: e}
}

func (s Service) Create(ctx *gofr.Context, e *models.Employee) (*models.Employee, error) {
	if e == nil {
		return nil, nil
	}

	err := e.Validate()
	if err != nil {
		return nil, err
	}

	emp, err := s.store.Get(ctx, *e.ID)
	if err != nil {
		return nil, err
	}

	if emp != nil {
		return e, nil
	}

	return s.store.Create(ctx, e)
}

func (s Service) Get(ctx *gofr.Context, id string) (*models.Employee, error) {
	if id == "" {
		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	empId, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	e, err := s.store.Get(ctx, int64(empId))
	if err != nil {
		return nil, err
	}

	if e == nil {
		return nil, http.ErrorEntityNotFound{Name: "id", Value: id}
	}

	return e, nil
}

func (s Service) GetAll(ctx *gofr.Context, page *pagination.Page) ([]models.Employee, error) {
	return s.store.GetAll(ctx, page)
}

func (s Service) Update(ctx *gofr.Context, e *models.Employee) (*models.Employee, error) {
	if e == nil {
		return nil, nil
	}

	err := e.Validate()
	if err != nil {
		return nil, err
	}

	emp, err := s.store.Get(ctx, *e.ID)
	if err != nil {
		return nil, err
	}

	if emp == nil {
		return nil, http.ErrorEntityNotFound{Name: "id", Value: fmt.Sprintf("%d", *e.ID)}
	}

	return s.store.Update(ctx, e)
}

func (s Service) Delete(ctx *gofr.Context, id string) error {
	if id == "" {
		return http.ErrorInvalidParam{Params: []string{"id"}}
	}

	empId, err := strconv.Atoi(id)
	if err != nil {
		return http.ErrorInvalidParam{Params: []string{"id"}}
	}

	return s.store.Delete(ctx, int64(empId))
}
