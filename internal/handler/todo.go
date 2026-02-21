package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sanketmote/go-todo-service/internal/service"
	spectodo "github.com/sanketmote/go-todo-service/internal/spec/todo"
	"github.com/sanketmote/go-todo-service/internal/svcerror"
	"github.com/sanketmote/gokit-wrapper/logger/svclog"
)

// TodoHandler handles HTTP for the Todo API.
type TodoHandler struct {
	svc    service.TodoService
	logger svclog.Logger
}

// NewTodoHandler returns a new TodoHandler.
func NewTodoHandler(svc service.TodoService, logger svclog.Logger) *TodoHandler {
	return &TodoHandler{svc: svc, logger: logger}
}

// writeJSON writes body as JSON with status code.
func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// writeError maps service errors to HTTP status and writes error body.
func (h *TodoHandler) writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, body := svcerror.HTTPStatusAndBody(err)
	if status >= 500 {
		h.logger.Error(r.Context(), "internal_error", err.Error(), "")
	}
	writeJSON(w, status, body)
}

// parseID extracts int64 id from mux URL param; returns (0, false) on invalid.
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

// Create handles POST /api/v1/todos.
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req spectodo.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, r, &svcerror.ValidationError{Message: "invalid JSON body"})
		return
	}
	resp, err := h.svc.Create(ctx, req.Title, req.Description, req.DueDate, req.RepeatType)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// List handles GET /api/v1/todos.
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	resp, err := h.svc.List(ctx, includeCompleted, sort, limit, offset)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetByID handles GET /api/v1/todos/:id.
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := parseID(r)
	if !ok {
		h.writeError(w, r, &svcerror.ValidationError{Message: "invalid id"})
		return
	}
	resp, err := h.svc.GetByID(ctx, id)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// Update handles PUT /api/v1/todos/:id.
func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := parseID(r)
	if !ok {
		h.writeError(w, r, &svcerror.ValidationError{Message: "invalid id"})
		return
	}
	var req spectodo.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, r, &svcerror.ValidationError{Message: "invalid JSON body"})
		return
	}
	resp, err := h.svc.Update(ctx, id, req.Title, req.Description, req.Completed, req.DueDate, req.RepeatType)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// Delete handles DELETE /api/v1/todos/:id.
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, ok := parseID(r)
	if !ok {
		h.writeError(w, r, &svcerror.ValidationError{Message: "invalid id"})
		return
	}
	err := h.svc.Delete(ctx, id)
	if err != nil {
		h.writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
