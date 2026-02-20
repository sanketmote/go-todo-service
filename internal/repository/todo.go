package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sanketmote/go-todo-service/internal/spec/datalayer/todo"
)

// TodoRepository defines data access for todos.
type TodoRepository interface {
	Create(ctx context.Context, title, description string, dueDate *time.Time, repeatType string) (*todo.Todo, error)
	GetByID(ctx context.Context, id int64) (*todo.Todo, error)
	List(ctx context.Context, includeCompleted bool, sort string, limit, offset int) ([]*todo.Todo, int, error)
	Update(ctx context.Context, id int64, title, description *string, completed *bool, dueDate *time.Time, repeatType *string) (*todo.Todo, error)
	Delete(ctx context.Context, id int64) error
}

type todoRepository struct {
	db *sqlx.DB
}

// NewTodoRepository returns a new TodoRepository.
func NewTodoRepository(db *sqlx.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, title, description string, dueDate *time.Time, repeatType string) (*todo.Todo, error) {
	if repeatType == "" {
		repeatType = "none"
	}
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO todos (title, description, completed, due_date, repeat_type) VALUES (?, ?, 0, ?, ?)`,
		title, description, dueDate, repeatType)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	var t todo.Todo
	err = r.db.GetContext(ctx, &t, `SELECT id, title, description, completed, due_date, repeat_type, created_at, updated_at FROM todos WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *todoRepository) GetByID(ctx context.Context, id int64) (*todo.Todo, error) {
	// Story 3.3
	return nil, nil
}

func (r *todoRepository) List(ctx context.Context, includeCompleted bool, sort string, limit, offset int) ([]*todo.Todo, int, error) {
	// Story 3.4
	return nil, 0, nil
}

func (r *todoRepository) Update(ctx context.Context, id int64, title, description *string, completed *bool, dueDate *time.Time, repeatType *string) (*todo.Todo, error) {
	// Story 3.5
	return nil, nil
}

func (r *todoRepository) Delete(ctx context.Context, id int64) error {
	// Story 3.6
	return nil
}
