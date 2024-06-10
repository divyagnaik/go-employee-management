package service

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"personal/go-employee-management/models"
	"personal/go-employee-management/pagination"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
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

func Test_handler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStr := NewMockemployeeStore(ctrl)
	mockSvc := New(mockStr)
	c, _ := container.NewMockContainer(t)

	id := int64(1)
	name := "test"
	req1 := models.Employee{
		ID:   &id,
		Name: &name,
	}

	req2 := models.Employee{
		Name: &name,
	}

	tests := []struct {
		name      string
		input     *models.Employee
		mockCalls []*gomock.Call
		want      *models.Employee
		wantErr   error
	}{
		{
			name:  "success case",
			input: &req1,
			mockCalls: []*gomock.Call{
				mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil),
				mockStr.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&req1, nil),
			},
			want: &req1,
		},
		{
			name: "success case - nil input",
		},
		{
			name:    "faiilure case - validation error",
			input:   &req2,
			wantErr: http.ErrorMissingParam{Params: []string{"employee.id"}},
		},
		{
			name:  "failure case - entity already exists",
			input: &req1,
			mockCalls: []*gomock.Call{
				mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&req1, nil),
			},
			want: &req1,
		},
		{
			name:  "failure case - error from GET",
			input: &req1,
			mockCalls: []*gomock.Call{
				mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error")),
			},
			wantErr: errors.New("some error"),
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("POST", "/test", "", nil, c)
			got, err := mockSvc.Create(ctx, tt.input)

			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)

		})
	}
}

func Test_handler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStr := NewMockemployeeStore(ctrl)
	mockSvc := New(mockStr)
	c, _ := container.NewMockContainer(t)

	id := int64(1)
	name := "test"
	req1 := models.Employee{
		ID:   &id,
		Name: &name,
	}

	tests := []struct {
		name      string
		id        string
		mockCalls *gomock.Call
		want      *models.Employee
		wantErr   error
	}{
		{
			name:      "success case",
			id:        "1",
			mockCalls: mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&req1, nil),
			want:      &req1,
		},
		{
			name:    "failure case - missing id",
			wantErr: http.ErrorInvalidParam(http.ErrorInvalidParam{Params: []string{"id"}}),
		},
		{
			name:    "failure case - invalid id",
			id:      "invalid",
			wantErr: http.ErrorInvalidParam(http.ErrorInvalidParam{Params: []string{"id"}}),
		},
		{
			name:      "error case - entity not found",
			id:        "1",
			mockCalls: mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil),
			wantErr:   http.ErrorEntityNotFound(http.ErrorEntityNotFound{Name: "id", Value: "1"}),
		},
		{
			name:      "error case - error from store",
			id:        "1",
			mockCalls: mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error")),
			wantErr:   errors.New("some error"),
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("GET", "/test", "1", nil, c)
			got, err := mockSvc.Get(ctx, tt.id)
			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func Test_handler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStr := NewMockemployeeStore(ctrl)
	mockService := New(mockStr)
	c, _ := container.NewMockContainer(t)
	ctx := createTestContext("GET", "/test", "", nil, c)

	id := int64(1)
	name := "test"
	req1 := []models.Employee{{
		ID:   &id,
		Name: &name,
	}}

	tests := []struct {
		name      string
		ctx       *gofr.Context
		mockCalls *gomock.Call
		want      []models.Employee
		wantErr   error
	}{
		{
			name:      "success case",
			ctx:       ctx,
			mockCalls: mockStr.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(req1, nil),
			want:      req1,
		},
		{
			name:      "error case - entity not found",
			mockCalls: mockStr.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error")),
			wantErr:   errors.New("some error"),
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockService.GetAll(ctx, &pagination.Page{Offset: 0, Size: 20})
			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func Test_handler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStr := NewMockemployeeStore(ctrl)
	mockSvc := New(mockStr)
	c, _ := container.NewMockContainer(t)

	id := int64(1)
	name := "test"
	req1 := models.Employee{
		ID:   &id,
		Name: &name,
	}

	req2 := models.Employee{
		Name: &name,
	}

	tests := []struct {
		name      string
		input     *models.Employee
		mockCalls []*gomock.Call
		want      *models.Employee
		wantErr   error
	}{
		{
			name:  "success case",
			input: &req1,
			mockCalls: []*gomock.Call{
				mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(&req1, nil),
				mockStr.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&req1, nil),
			},
			want: &req1,
		},
		{
			name: "success case - nil input",
		},
		{
			name:    "faiilure case - validation error",
			input:   &req2,
			wantErr: http.ErrorMissingParam{Params: []string{"employee.id"}},
		},
		{
			name:  "failure case - entity does not exists",
			input: &req1,
			mockCalls: []*gomock.Call{
				mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil),
			},
			wantErr: http.ErrorEntityNotFound(http.ErrorEntityNotFound{Name: "id", Value: "1"}),
		},
		{
			name:  "failure case - error from GET",
			input: &req1,
			mockCalls: []*gomock.Call{
				mockStr.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error")),
			},
			wantErr: errors.New("some error"),
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("POST", "/test", "", nil, c)
			got, err := mockSvc.Update(ctx, tt.input)

			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)

		})
	}
}

func Test_handler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockStr := NewMockemployeeStore(ctrl)
	mockSvc := New(mockStr)
	c, _ := container.NewMockContainer(t)

	tests := []struct {
		name      string
		id        string
		mockCalls *gomock.Call
		want      interface{}
		wantErr   error
	}{
		{
			name:      "success case",
			id:        "1",
			mockCalls: mockStr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil),
		},
		{
			name:    "failure case - missing id",
			wantErr: http.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:    "failure case - invalid id",
			id:      "invalid",
			wantErr: http.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:      "error case - error from store",
			id:        "1",
			mockCalls: mockStr.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("some error")),
			wantErr:   errors.New("some error"),
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("DELETE", "/test", "1", nil, c)
			err := mockSvc.Delete(ctx, tt.id)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}
