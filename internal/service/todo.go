package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/sanketmote/go-todo-service/internal/repository"
	"github.com/sanketmote/go-todo-service/internal/spec/datalayer/todo"
	spectodo "github.com/sanketmote/go-todo-service/internal/spec/todo"
	"github.com/sanketmote/go-todo-service/internal/svcerror"
	"github.com/sanketmote/gokit-wrapper/logger/svclog"
)

// TodoService defines business logic for todos.
type TodoService interface {
	Create(ctx context.Context, title, description string, dueDate *time.Time, repeatType string) (*spectodo.TodoResponse, error)
	GetByID(ctx context.Context, id int64) (*spectodo.TodoResponse, error)
	List(ctx context.Context, includeCompleted bool, sort string, limit, offset int) (*spectodo.ListTodosResponse, error)
	Update(ctx context.Context, id int64, title, description *string, completed *bool, dueDate *time.Time, repeatType *string) (*spectodo.TodoResponse, error)
	Delete(ctx context.Context, id int64) error
}

type todoService struct {
	repo   repository.TodoRepository
	logger svclog.Logger
}

// NewTodoService returns a new TodoService.
func NewTodoService(repo repository.TodoRepository, logger svclog.Logger) TodoService {
	return &todoService{repo: repo, logger: logger}
}

func (s *todoService) Create(ctx context.Context, title, description string, dueDate *time.Time, repeatType string) (*spectodo.TodoResponse, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, &svcerror.ValidationError{Message: "title is required"}
	}
	if len(title) > 500 {
		return nil, &svcerror.ValidationError{Message: "title must be 1-500 characters"}
	}
	if !validRepeatType(repeatType) {
		return nil, &svcerror.ValidationError{Message: "repeat_type must be one of: none, daily, weekly, monthly, yearly"}
	}
	t, err := s.repo.Create(ctx, title, description, dueDate, repeatType)
	if err != nil {
		return nil, err
	}
	return toTodoResponse(t), nil
}

func (s *todoService) GetByID(ctx context.Context, id int64) (*spectodo.TodoResponse, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, svcerror.ErrTodoNotFound) {
			return nil, svcerror.ErrTodoNotFound
		}
		return nil, err
	}
	return toTodoResponse(t), nil
}

func (s *todoService) List(ctx context.Context, includeCompleted bool, sort string, limit, offset int) (*spectodo.ListTodosResponse, error) {
	if sort == "" {
		sort = spectodo.SortDueDateAsc
	}
	if !validSort(sort) {
		return nil, &svcerror.ValidationError{Message: "sort must be one of: due_date_asc, due_date_desc, created_at_asc, created_at_desc"}
	}
	if limit <= 0 {
		limit = spectodo.LimitDefault
	}
	if limit > spectodo.LimitMax {
		limit = spectodo.LimitMax
	}
	if offset < 0 {
		offset = 0
	}
	items, total, err := s.repo.List(ctx, includeCompleted, sort, limit, offset)
	if err != nil {
		return nil, err
	}
	resp := &spectodo.ListTodosResponse{
		Items: make([]spectodo.TodoResponse, 0, len(items)),
		Total: int64(total),
	}
	for _, t := range items {
		resp.Items = append(resp.Items, *toTodoResponse(t))
	}
	return resp, nil
}

func (s *todoService) Update(ctx context.Context, id int64, title, description *string, completed *bool, dueDate *time.Time, repeatType *string) (*spectodo.TodoResponse, error) {
	if title != nil {
		t := strings.TrimSpace(*title)
		title = &t
		if *title == "" {
			return nil, &svcerror.ValidationError{Message: "title must be 1-500 characters"}
		}
		if len(*title) > 500 {
			return nil, &svcerror.ValidationError{Message: "title must be 1-500 characters"}
		}
	}
	if repeatType != nil && !validRepeatType(*repeatType) {
		return nil, &svcerror.ValidationError{Message: "repeat_type must be one of: none, daily, weekly, monthly, yearly"}
	}
	t, err := s.repo.Update(ctx, id, title, description, completed, dueDate, repeatType)
	if err != nil {
		if errors.Is(err, svcerror.ErrTodoNotFound) {
			return nil, svcerror.ErrTodoNotFound
		}
		return nil, err
	}
	return toTodoResponse(t), nil
}

func (s *todoService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, svcerror.ErrTodoNotFound) {
			return svcerror.ErrTodoNotFound
		}
		return err
	}
	return nil
}

func validRepeatType(r string) bool {
	switch r {
	case "", spectodo.RepeatNone, spectodo.RepeatDaily, spectodo.RepeatWeekly, spectodo.RepeatMonthly, spectodo.RepeatYearly:
		return true
	default:
		return false
	}
}

func validSort(sort string) bool {
	switch sort {
	case spectodo.SortDueDateAsc, spectodo.SortDueDateDesc, spectodo.SortCreatedAsc, spectodo.SortCreatedDesc:
		return true
	default:
		return false
	}
}

func toTodoResponse(t *todo.Todo) *spectodo.TodoResponse {
	if t == nil {
		return nil
	}
	return &spectodo.TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		DueDate:     t.DueDate,
		RepeatType:  t.RepeatType,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
