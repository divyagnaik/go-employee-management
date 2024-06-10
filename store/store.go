package store

import (
	"database/sql"
	"personal/go-employee-management/models"
	"personal/go-employee-management/pagination"

	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/datasource"
)

type Store struct{}

func New() *Store {
	return &Store{}
}

func (s Store) Create(ctx *gofr.Context, e *models.Employee) (*models.Employee, error) {
	if e == nil {
		return nil, nil
	}

	_, err := ctx.SQL.ExecContext(ctx, "INSERT INTO employee (id, name, position, salary) values ($1, $2, $3, $4);", e.ID, e.Name, e.Position, e.Salary)
	if err != nil {
		return nil, datasource.ErrorDB{Message: "Internal Server Error", Err: err}
	}

	return e, nil
}

func (s Store) Get(ctx *gofr.Context, id int64) (*models.Employee, error) {
	row := ctx.SQL.QueryRowContext(ctx, "SELECT id, name, position, salary FROM employee where id = $1;", id)

	var e models.Employee

	err := row.Scan(&e.ID, &e.Name, &e.Position, &e.Salary)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, datasource.ErrorDB{Message: "Internal Server Error", Err: err}
	}

	return &e, nil
}

func (s Store) GetAll(ctx *gofr.Context, page *pagination.Page) ([]models.Employee, error) {
	rows, err := ctx.SQL.QueryContext(ctx, "SELECT id, name, position, salary FROM employee OFFSET $1 LIMIT $2;", page.Offset, page.Size)
	if err != nil {
		return nil, datasource.ErrorDB{Message: "Internal Server Error", Err: err}
	}

	defer rows.Close()

	var resp []models.Employee

	for rows.Next() {
		var e models.Employee
		err = rows.Scan(&e.ID, &e.Name, &e.Position, &e.Salary)
		if err != nil {
			return nil, datasource.ErrorDB{Message: "Internal Server Error", Err: err}
		}

		resp = append(resp, e)
	}

	return resp, nil
}
