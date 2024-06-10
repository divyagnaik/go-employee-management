package models

import (
	"personal/go-employee-management/pagination"

	"gofr.dev/pkg/gofr/http"
)

type Employee struct {
	ID       *int64  `json:"id"`
	Name     *string `json:"name"`
	Position string  `json:"position"`
	Salary   int64   `json:"salary"`
}

func (e *Employee) Validate() error {
	if e == nil {
		return http.ErrorMissingParam{Params: []string{"employee"}}
	}

	if e.ID == nil {
		return http.ErrorMissingParam{Params: []string{"employee.id"}}
	}

	if e.Name == nil {
		return http.ErrorMissingParam{Params: []string{"employee.name"}}
	}

	if *e.Name == "" {
		return http.ErrorInvalidParam{Params: []string{"employee.name"}}
	}

	return nil
}

type GetAllResp struct {
	Meta pagination.Meta `json:"meta"`
	Resp []Employee      `json:"data"`
}
