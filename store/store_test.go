package store

import (
	"bytes"
	"database/sql"
	"net/http/httptest"
	"personal/go-employee-management/models"
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
	"gofr.dev/pkg/gofr/logging"
)

type DB struct {
	// contains unexported or private fields
	*sql.DB
	logger  datasource.Logger
	config  *gofrsql.DBConfig
	metrics gofrsql.Metrics
}

func getDB(t *testing.T, logLevel logging.Level) (*DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual), sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db := &DB{mockDB, logging.NewMockLogger(logLevel), nil, nil}
	db.config = &gofrsql.DBConfig{}

	return db, mock
}

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
	db, mock := getDB(t, logging.INFO)
	defer db.DB.Close()
	ctrl := gomock.NewController(t)
	c, _ := container.NewMockContainer(t)
	c.SQL = container.NewMockDB(ctrl)

	ctx := createTestContext("POST", "", "", nil, c)

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
		mockCall *sqlmock.ExpectedExec
		want     *models.Employee
		wantErr  error
	}{
		{
			name:  "success",
			input: &req1,
			mockCall: mock.ExpectExec("INSERT INTO employee").
				WithArgs(1, "John", "Developer", 50000).
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Store{}
			got, err := s.Create(ctx, tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
