package todo

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	spectodo "github.com/sanketmote/go-todo-service/internal/spec/todo"
	"github.com/sanketmote/go-todo-service/internal/svcerror"
)

// Path constants for router registration.
const (
	PathTodos     = "/api/v1/todos"
	PathTodosByID = "/api/v1/todos/{id}"
)

// Request types for decode (path/query params combined with body where needed).
type listRequest struct {
	IncludeCompleted bool
	Sort             string
	Limit            int
	Offset           int
}

type getByIDRequest struct {
	ID int64
}

type updateRequest struct {
	ID          int64
	Title       *string
	Description *string
	Completed   *bool
	DueDate     *time.Time
	RepeatType  *string
}

type deleteRequest struct {
	ID int64
}

// AddHandlers registers Todo API routes: path → decode → endpoint → encode.
func AddHandlers(r *mux.Router, endpoints spectodo.Endpoints) {
	r.Methods("POST").Path(PathTodos).Handler(httptransport.NewServer(
		endpoints.Create,
		decodeCreateRequest,
		encodeTodoResponse(http.StatusCreated),
		httptransport.ServerErrorEncoder(encodeError),
	))
	r.Methods("GET").Path(PathTodos).Handler(httptransport.NewServer(
		endpoints.List,
		decodeListRequest,
		encodeListResponse,
		httptransport.ServerErrorEncoder(encodeError),
	))
	r.Methods("GET").Path(PathTodosByID).Handler(httptransport.NewServer(
		endpoints.GetByID,
		decodeGetByIDRequest,
		encodeTodoResponse(http.StatusOK),
		httptransport.ServerErrorEncoder(encodeError),
	))
	r.Methods("PUT").Path(PathTodosByID).Handler(httptransport.NewServer(
		endpoints.Update,
		decodeUpdateRequest,
		encodeTodoResponse(http.StatusOK),
		httptransport.ServerErrorEncoder(encodeError),
	))
	r.Methods("DELETE").Path(PathTodosByID).Handler(httptransport.NewServer(
		endpoints.Delete,
		decodeDeleteRequest,
		encodeNoContent,
		httptransport.ServerErrorEncoder(encodeError),
	))
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req spectodo.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, &svcerror.ValidationError{Message: "invalid JSON body"}
	}
	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	includeCompleted := q.Get(spectodo.ParamIncludeCompleted) == "true" || q.Get(spectodo.ParamIncludeCompleted) == "1"
	sort := q.Get(spectodo.ParamSort)
	if sort == "" {
		sort = spectodo.SortDueDateAsc
	}
	limit, _ := strconv.Atoi(q.Get(spectodo.ParamLimit))
	if limit <= 0 {
		limit = spectodo.LimitDefault
	}
	if limit > spectodo.LimitMax {
		limit = spectodo.LimitMax
	}
	offset, _ := strconv.Atoi(q.Get(spectodo.ParamOffset))
	if offset < 0 {
		offset = 0
	}
	return listRequest{IncludeCompleted: includeCompleted, Sort: sort, Limit: limit, Offset: offset}, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, ok := parseID(r)
	if !ok {
		return nil, &svcerror.ValidationError{Message: "invalid id"}
	}
	return getByIDRequest{ID: id}, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, ok := parseID(r)
	if !ok {
		return nil, &svcerror.ValidationError{Message: "invalid id"}
	}
	var body spectodo.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, &svcerror.ValidationError{Message: "invalid JSON body"}
	}
	return updateRequest{ID: id, Title: body.Title, Description: body.Description, Completed: body.Completed, DueDate: body.DueDate, RepeatType: body.RepeatType}, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, ok := parseID(r)
	if !ok {
		return nil, &svcerror.ValidationError{Message: "invalid id"}
	}
	return deleteRequest{ID: id}, nil
}

func parseID(r *http.Request) (int64, bool) {
	s := mux.Vars(r)["id"]
	if s == "" {
		return 0, false
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func encodeTodoResponse(status int) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		return json.NewEncoder(w).Encode(response)
	}
}

func encodeListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func encodeNoContent(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	status, body := svcerror.HTTPStatusAndBody(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
