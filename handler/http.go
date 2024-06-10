package handler

import (
	"personal/go-employee-management/models"
	"personal/go-employee-management/pagination"
	"strconv"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/http"
)

type handler struct {
	service employeeService
}

func New(e employeeService) *handler {
	return &handler{service: e}
}

func (h handler) Create(ctx *gofr.Context) (interface{}, error) {
	var req models.Employee
	err := ctx.Bind(&req)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"Request"}}
	}

	resp, err := h.service.Create(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h handler) Get(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	resp, err := h.service.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h handler) GetAll(ctx *gofr.Context) (interface{}, error) {
	page, err := pagination.Pagination(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := h.service.GetAll(ctx, page)
	if err != nil {
		return nil, err
	}

	return models.GetAllResp{Meta: pagination.Meta{Page: *page}, Resp: resp}, nil
}

func (h handler) Update(ctx *gofr.Context) (interface{}, error) {
	var req models.Employee

	err := ctx.Bind(&req)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"Request"}}
	}

	id := ctx.PathParam("id")

	employeeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.ErrorInvalidParam{Params: []string{"id"}}
	}

	employeeIDInt64 := int64(employeeID)
	req.ID = &employeeIDInt64

	resp, err := h.service.Update(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h handler) Delete(ctx *gofr.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	return nil, h.service.Delete(ctx, id)
}
