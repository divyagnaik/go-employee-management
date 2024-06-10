package main

import (
	"personal/go-employee-management/handler"
	"personal/go-employee-management/migrations"
	"personal/go-employee-management/service"
	"personal/go-employee-management/store"

	"gofr.dev/pkg/gofr"
)

func main() {
	g := gofr.New()

	g.Migrate(migrations.All())

	str := store.New()
	svc := service.New(str)
	h := handler.New(svc)

	g.POST("/empoylee/v1/empolyee", h.Create)
	g.GET("/empoylee/v1/empolyee/{id}", h.Get)
	g.GET("/empoylee/v1/empolyee", h.GetAll)
	g.PUT("/employee/v1/employee/{id}", h.Update)
	g.DELETE("/employee/v1/employee/{id}", h.Delete)

	g.Run()
}
