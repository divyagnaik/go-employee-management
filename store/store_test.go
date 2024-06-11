package store

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http/httptest"
	"personal/go-employee-management/models"
	"personal/go-employee-management/pagination"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	"gofr.dev/pkg/gofr/datasource"
	gofrsql "gofr.dev/pkg/gofr/datasource/sql"
	"gofr.dev/pkg/gofr/http"
)

func createTestContext(method, path, id string, body []byte, cont *container.Container) *gofr.Context {
	testReq := httptest.NewRequest(method, path+"/"+id, bytes.NewBuffer(body))
	testReq = mux.SetURLVars(testReq, map[string]string{"id": id})
	testReq.Header.Set("Content-Type", "application/json")
	gofrReq := http.NewRequest(testReq)
	return &gofr.Context{
		Context:   gofrReq.Context(),
		Request:   gofrReq,
		Container: cont,
	}
}
func TestStore_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c, _ := container.NewMockContainer(t)
	ctx := createTestContext("POST", "", "", nil, c)

	db, mock, mockMetric := gofrsql.NewSQLMocksWithConfig(t, &gofrsql.DBConfig{})
	defer db.Close()

	ctx.SQL = db

	id := int64(1)
	name := "John"

	req1 := models.Employee{
		ID:       &id,
		Name:     &name,
		Salary:   50000,
		Position: "Developer",
	}

	tests := []struct {
		name     string
		input    *models.Employee
		mockCall []interface{}
		want     *models.Employee
		wantErr  error
	}{
		{
			name:  "success case",
			input: &req1,
			mockCall: []interface{}{
				mock.ExpectExec("INSERT INTO employee (id, name, position, salary) values ($1, $2, $3, $4);").WithArgs(1, "John", "Developer", 50000).
					WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			want: &req1,
		},
		{
			name:  "failure case",
			input: &req1,
			mockCall: []interface{}{
				mock.ExpectExec("INSERT INTO employee (id, name, position, salary) values ($1, $2, $3, $4);").WithArgs(1, "John", "Developer", 50000).
					WillReturnError(errors.New("some error")),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			wantErr: datasource.ErrorDB{Message: "Internal Server Error", Err: errors.New("some error")},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Store{}
			got, err := s.Create(ctx, tt.input)

			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func TestStore_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c, _ := container.NewMockContainer(t)
	ctx := createTestContext("GET", "", "1", nil, c)

	db, mock, mockMetric := gofrsql.NewSQLMocksWithConfig(t, &gofrsql.DBConfig{})
	defer db.Close()

	ctx.SQL = db

	id := int64(1)
	name := "John"
	position := "Developer"

	tests := []struct {
		name     string
		id       int64
		mockCall []interface{}
		want     *models.Employee
		wantErr  error
	}{
		{
			name: "success case",
			id:   id,
			mockCall: []interface{}{
				mock.ExpectQuery("SELECT id, name, position, salary FROM employee where id = $1;").WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary"}).AddRow(id, name, position, 50000)),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			want: &models.Employee{
				ID:       &id,
				Name:     &name,
				Position: position,
				Salary:   50000,
			},
		},
		{
			name: "no rows case",
			id:   id,
			mockCall: []interface{}{
				mock.ExpectQuery("SELECT id, name, position, salary FROM employee where id = $1;").WithArgs(id).
					WillReturnError(sql.ErrNoRows),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "failure case",
			id:   id,
			mockCall: []interface{}{
				mock.ExpectQuery("SELECT id, name, position, salary FROM employee where id = $1;").WithArgs(id).
					WillReturnError(errors.New("some error")),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			want:    nil,
			wantErr: datasource.ErrorDB{Message: "Internal Server Error", Err: errors.New("some error")},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Store{}
			got, err := s.Get(ctx, tt.id)

			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func TestStore_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c, _ := container.NewMockContainer(t)
	ctx := createTestContext("GET", "", "", nil, c)

	db, mock, mockMetric := gofrsql.NewSQLMocksWithConfig(t, &gofrsql.DBConfig{})
	defer db.Close()

	ctx.SQL = db

	page := &pagination.Page{
		Offset: 0,
		Size:   2,
	}

	id1, id2 := int64(1), int64(2)
	name1, name2 := "John", "Doe"
	position1, position2 := "Developer", "Manager"

	employees := []models.Employee{
		{
			ID:       &id1,
			Name:     &name1,
			Position: position1,
			Salary:   50000,
		},
		{
			ID:       &id2,
			Name:     &name2,
			Position: position2,
			Salary:   60000,
		},
	}

	tests := []struct {
		name     string
		page     *pagination.Page
		mockCall []interface{}
		want     []models.Employee
		wantErr  error
	}{
		{
			name: "success case",
			page: page,
			mockCall: []interface{}{
				mock.ExpectQuery("SELECT id, name, position, salary FROM employee OFFSET $1 LIMIT $2;").
					WithArgs(page.Offset, page.Size).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary"}).
						AddRow(id1, name1, position1, 50000).
						AddRow(id2, name2, position2, 60000)),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			want: employees,
		},
		{
			name: "failure case - scan error",
			page: page,
			mockCall: []interface{}{
				mock.ExpectQuery("SELECT id, name, position, salary FROM employee OFFSET $1 LIMIT $2;").
					WithArgs(page.Offset, page.Size).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary"}).
						AddRow("invalid_id", name1, position1, 50000). // This row will cause a scan error
						AddRow(id2, name2, position2, 60000)),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			want:    nil,
			wantErr: errors.New("Internal Server Error: sql: Scan error on column index 0, name \"id\": converting driver.Value type string (\"invalid_id\") to a int64: invalid syntax")},
		{
			name: "failure case",
			page: page,
			mockCall: []interface{}{
				mock.ExpectQuery("SELECT id, name, position, salary FROM employee OFFSET $1 LIMIT $2;").
					WithArgs(page.Offset, page.Size).
					WillReturnError(errors.New("some error")),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			want:    nil,
			wantErr: datasource.ErrorDB{Message: "Internal Server Error", Err: errors.New("some error")},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Store{}
			got, err := s.GetAll(ctx, tt.page)

			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			if i == 1 {
				assert.Equal(t, tt.wantErr.Error(), err.Error(), "TEST[%d], Failed.\n%s", i, tt.name)
				return
			}
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func TestStore_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c, _ := container.NewMockContainer(t)
	ctx := createTestContext("PUT", "", "", nil, c)

	db, mock, mockMetric := gofrsql.NewSQLMocksWithConfig(t, &gofrsql.DBConfig{})
	defer db.Close()

	ctx.SQL = db

	id := int64(1)
	name := "John"
	position := "Senior Developer"

	employee := models.Employee{
		ID:       &id,
		Name:     &name,
		Position: position,
		Salary:   70000,
	}

	tests := []struct {
		name     string
		input    *models.Employee
		mockCall []interface{}
		want     *models.Employee
		wantErr  error
	}{
		{
			name:  "success case",
			input: &employee,
			mockCall: []interface{}{
				mock.ExpectExec("UPDATE employee SET name=$1, position=$2, salary=$3 WHERE id=$4;").
					WithArgs(name, position, 70000, id).
					WillReturnResult(sqlmock.NewResult(0, 1)),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			want: &employee,
		},
		{
			name:  "failure case",
			input: &employee,
			mockCall: []interface{}{
				mock.ExpectExec("UPDATE employee SET name=$1, position=$2, salary=$3 WHERE id=$4;").
					WithArgs(name, position, 70000, id).
					WillReturnError(errors.New("some error")),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			wantErr: datasource.ErrorDB{Message: "Internal Server Error", Err: errors.New("some error")},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Store{}
			got, err := s.Update(ctx, tt.input)

			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func TestStore_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c, _ := container.NewMockContainer(t)
	ctx := createTestContext("DELETE", "", "1", nil, c)

	db, mock, mockMetric := gofrsql.NewSQLMocksWithConfig(t, &gofrsql.DBConfig{})
	defer db.Close()

	ctx.SQL = db

	id := int64(1)

	tests := []struct {
		name     string
		id       int64
		mockCall []interface{}
		wantErr  error
	}{
		{
			name: "success case",
			id:   id,
			mockCall: []interface{}{
				mock.ExpectExec("DELETE FROM employee where id = $1;").
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1)),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			wantErr: nil,
		},
		{
			name: "failure case",
			id:   id,
			mockCall: []interface{}{
				mock.ExpectExec("DELETE FROM employee where id = $1;").
					WithArgs(id).
					WillReturnError(errors.New("some error")),
				mockMetric.EXPECT().RecordHistogram(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
			},
			wantErr: datasource.ErrorDB{Message: "Internal Server Error", Err: errors.New("some error")},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Store{}
			err := s.Delete(ctx, tt.id)

			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}
