package handler

import (
	"bytes"
	"encoding/json"
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
	mockSvc := NewMockemployeeService(ctrl)
	mockHandler := New(mockSvc)
	c, _ := container.NewMockContainer(t)

	id := int64(1)
	name := "test"
	req1 := models.Employee{
		ID:   &id,
		Name: &name,
	}

	reqBody1, _ := json.Marshal(req1)

	reqBody2 := []byte("test")

	tests := []struct {
		name      string
		reqBody   []byte
		mockCalls *gomock.Call
		want      interface{}
		wantErr   error
	}{
		{
			name:      "success case",
			reqBody:   reqBody1,
			mockCalls: mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&req1, nil),
			want:      &req1,
		},
		{
			name:    "faiilure case - bind error",
			reqBody: reqBody2,
			wantErr: http.ErrorInvalidParam{Params: []string{"Request"}},
		},
		{
			name:      "failure case - entity already exists",
			reqBody:   reqBody1,
			mockCalls: mockSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, http.ErrorEntityAlreadyExist{}),
			want:      nil,
			wantErr:   http.ErrorEntityAlreadyExist{},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("POST", "/test", "", tt.reqBody, c)
			got, err := mockHandler.Create(ctx)

			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)

		})
	}
}

func Test_handler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := NewMockemployeeService(ctrl)
	mockHandler := New(mockSvc)
	c, _ := container.NewMockContainer(t)

	id := int64(1)
	name := "test"
	req1 := models.Employee{
		ID:   &id,
		Name: &name,
	}

	tests := []struct {
		name      string
		mockCalls *gomock.Call
		want      interface{}
		wantErr   error
	}{
		{
			name:      "success case",
			mockCalls: mockSvc.EXPECT().Get(gomock.Any(), "1").Return(&req1, nil),
			want:      &req1,
		},
		{
			name:      "error case - entity not found",
			mockCalls: mockSvc.EXPECT().Get(gomock.Any(), "1").Return(nil, http.ErrorEntityNotFound{}),
			wantErr:   http.ErrorEntityNotFound{},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("GET", "/test", "1", nil, c)
			got, err := mockHandler.Get(ctx)
			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func Test_handler_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := NewMockemployeeService(ctrl)
	mockHandler := New(mockSvc)
	c, _ := container.NewMockContainer(t)
	ctx := createTestContext("GET", "/test", "", nil, c)
	ctx1 := createTestContext("GET", "/test?page.size=@!", "", nil, c)

	id := int64(1)
	name := "test"
	req1 := []models.Employee{{
		ID:   &id,
		Name: &name,
	}}
	expResp := models.GetAllResp{
		Resp: req1,
		Meta: pagination.Meta{Page: pagination.Page{Offset: 0, Size: 20}},
	}

	tests := []struct {
		name      string
		ctx       *gofr.Context
		mockCalls *gomock.Call
		want      interface{}
		wantErr   error
	}{
		{
			name:      "success case",
			ctx:       ctx,
			mockCalls: mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(req1, nil),
			want:      expResp,
		},
		{
			name:      "error case - entity not found",
			ctx:       ctx,
			mockCalls: mockSvc.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, http.ErrorEntityNotFound{}),
			wantErr:   http.ErrorEntityNotFound{},
		},
		{
			name:    "error case - invalid pagination data",
			ctx:     ctx1,
			wantErr: http.ErrorInvalidParam(http.ErrorInvalidParam{Params: []string{"page.size"}}),
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockHandler.GetAll(tt.ctx)
			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func Test_handler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := NewMockemployeeService(ctrl)
	mockHandler := New(mockSvc)
	c, _ := container.NewMockContainer(t)

	id := int64(1)
	name := "test"
	req1 := models.Employee{
		ID:   &id,
		Name: &name,
	}

	reqBody1, _ := json.Marshal(req1)

	reqBody2 := []byte("test")

	tests := []struct {
		name      string
		id        string
		reqBody   []byte
		mockCalls *gomock.Call
		want      interface{}
		wantErr   error
	}{
		{
			name:      "success case",
			id:        "1",
			reqBody:   reqBody1,
			mockCalls: mockSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(&req1, nil),
			want:      &req1,
		},
		{
			name:    "faiilure case - bind error",
			id:      "1",
			reqBody: reqBody2,
			wantErr: http.ErrorInvalidParam{Params: []string{"Request"}},
		},
		{
			name:    "faiilure case - invalid id",
			id:      "invalid",
			reqBody: reqBody1,
			wantErr: http.ErrorInvalidParam{Params: []string{"id"}},
		},
		{
			name:      "failure case - entity does not exists",
			id:        "1",
			reqBody:   reqBody1,
			mockCalls: mockSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil, http.ErrorEntityNotFound{}),
			want:      nil,
			wantErr:   http.ErrorEntityNotFound{},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("PUT", "/test", tt.id, tt.reqBody, c)
			got, err := mockHandler.Update(ctx)
			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}

func Test_handler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := NewMockemployeeService(ctrl)
	mockHandler := New(mockSvc)
	c, _ := container.NewMockContainer(t)

	tests := []struct {
		name      string
		mockCalls *gomock.Call
		want      interface{}
		wantErr   error
	}{
		{
			name:      "success case",
			mockCalls: mockSvc.EXPECT().Delete(gomock.Any(), "1").Return(nil),
		},
		{
			name:      "error case - entity not found",
			mockCalls: mockSvc.EXPECT().Delete(gomock.Any(), "1").Return(http.ErrorEntityNotFound{}),
			wantErr:   http.ErrorEntityNotFound{},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := createTestContext("DELETE", "/test", "1", nil, c)
			got, err := mockHandler.Delete(ctx)
			assert.Equal(t, tt.want, got, "TEST[%d], Failed.\n%s", i, tt.name)
			assert.Equal(t, tt.wantErr, err, "TEST[%d], Failed.\n%s", i, tt.name)
		})
	}
}
