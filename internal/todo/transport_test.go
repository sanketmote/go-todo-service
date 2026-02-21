package todo

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	spectodo "github.com/sanketmote/go-todo-service/internal/spec/todo"
	"github.com/sanketmote/go-todo-service/internal/svcerror"
)

// mockHTTPService returns fixed responses for HTTP tests.
type mockHTTPService struct{}

func (m *mockHTTPService) Create(ctx context.Context, title, desc string, due *time.Time, repeat string) (*spectodo.TodoResponse, error) {
	if title == "" {
		return nil, &svcerror.ValidationError{Message: "title is required"}
	}
	return &spectodo.TodoResponse{ID: 1, Title: title, RepeatType: repeat, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockHTTPService) GetByID(ctx context.Context, id int64) (*spectodo.TodoResponse, error) {
	if id == 999 {
		return nil, svcerror.ErrTodoNotFound
	}
	return &spectodo.TodoResponse{ID: id, Title: "Task", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockHTTPService) List(ctx context.Context, includeCompleted bool, sort string, limit, offset int) (*spectodo.ListTodosResponse, error) {
	return &spectodo.ListTodosResponse{Items: []spectodo.TodoResponse{}, Total: 0}, nil
}

func (m *mockHTTPService) Update(ctx context.Context, id int64, title, desc *string, completed *bool, due *time.Time, repeat *string) (*spectodo.TodoResponse, error) {
	if id == 999 {
		return nil, svcerror.ErrTodoNotFound
	}
	t := "Updated"
	if title != nil {
		t = *title
	}
	return &spectodo.TodoResponse{ID: id, Title: t, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
}

func (m *mockHTTPService) Delete(ctx context.Context, id int64) error {
	if id == 999 {
		return svcerror.ErrTodoNotFound
	}
	return nil
}

func TestTransport_HTTPHandlers(t *testing.T) {
	svc := &mockHTTPService{}
	eps := MakeEndpoints(svc)
	r := mux.NewRouter()
	AddHandlers(r, eps)

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantBody   func(t *testing.T, body []byte)
	}{
		{
			name:       "POST create ok",
			method:     "POST",
			path:       "/api/v1/todos",
			body:       `{"title":"Buy groceries","repeat_type":"none"}`,
			wantStatus: http.StatusCreated,
			wantBody: func(t *testing.T, body []byte) {
				var resp spectodo.TodoResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if resp.ID != 1 || resp.Title != "Buy groceries" {
					t.Errorf("got %+v", resp)
				}
			},
		},
		{
			name:       "POST create invalid JSON",
			method:     "POST",
			path:       "/api/v1/todos",
			body:       `{invalid`,
			wantStatus: http.StatusBadRequest,
			wantBody: func(t *testing.T, body []byte) {
				var m map[string]string
				if err := json.Unmarshal(body, &m); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if m["error"] != "validation_error" {
					t.Errorf("got %v", m)
				}
			},
		},
		{
			name:       "GET list ok",
			method:     "GET",
			path:       "/api/v1/todos",
			wantStatus: http.StatusOK,
			wantBody: func(t *testing.T, body []byte) {
				var resp spectodo.ListTodosResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if resp.Total != 0 {
					t.Errorf("Total = %d, want 0", resp.Total)
				}
			},
		},
		{
			name:       "GET list with query params",
			method:     "GET",
			path:       "/api/v1/todos?include_completed=true&sort=due_date_desc&limit=10&offset=0",
			wantStatus: http.StatusOK,
			wantBody:   nil,
		},
		{
			name:       "GET by id ok",
			method:     "GET",
			path:       "/api/v1/todos/1",
			wantStatus: http.StatusOK,
			wantBody: func(t *testing.T, body []byte) {
				var resp spectodo.TodoResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if resp.ID != 1 {
					t.Errorf("ID = %d, want 1", resp.ID)
				}
			},
		},
		{
			name:       "GET by id not found",
			method:     "GET",
			path:       "/api/v1/todos/999",
			wantStatus: http.StatusNotFound,
			wantBody: func(t *testing.T, body []byte) {
				var m map[string]string
				if err := json.Unmarshal(body, &m); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if m["error"] != "not_found" {
					t.Errorf("got %v", m)
				}
			},
		},
		{
			name:       "GET by id invalid",
			method:     "GET",
			path:       "/api/v1/todos/abc",
			wantStatus: http.StatusBadRequest,
			wantBody: func(t *testing.T, body []byte) {
				var m map[string]string
				if err := json.Unmarshal(body, &m); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if m["error"] != "validation_error" {
					t.Errorf("got %v", m)
				}
			},
		},
		{
			name:       "PUT update ok",
			method:     "PUT",
			path:       "/api/v1/todos/1",
			body:       `{"title":"Updated","completed":true}`,
			wantStatus: http.StatusOK,
			wantBody: func(t *testing.T, body []byte) {
				var resp spectodo.TodoResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("unmarshal: %v", err)
				}
				if resp.Title != "Updated" {
					t.Errorf("Title = %q, want Updated", resp.Title)
				}
			},
		},
		{
			name:       "PUT update not found",
			method:     "PUT",
			path:       "/api/v1/todos/999",
			body:       `{"title":"Updated"}`,
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
		{
			name:       "DELETE ok",
			method:     "DELETE",
			path:       "/api/v1/todos/1",
			wantStatus: http.StatusNoContent,
			wantBody:   nil,
		},
		{
			name:       "DELETE not found",
			method:     "DELETE",
			path:       "/api/v1/todos/999",
			wantStatus: http.StatusNotFound,
			wantBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error
			if tt.body != "" {
				req, err = http.NewRequest(tt.method, "http://test"+tt.path, bytes.NewReader([]byte(tt.body)))
			} else {
				req, err = http.NewRequest(tt.method, "http://test"+tt.path, nil)
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d\nbody: %s", rec.Code, tt.wantStatus, rec.Body.Bytes())
			}
			if tt.wantBody != nil && rec.Body.Len() > 0 {
				tt.wantBody(t, rec.Body.Bytes())
			}
		})
	}
}
