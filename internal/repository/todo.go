package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sanketmote/go-todo-service/internal/spec/datalayer/todo"
	"github.com/sanketmote/go-todo-service/internal/svcerror"
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
	var t todo.Todo
	err := r.db.GetContext(ctx, &t, `SELECT id, title, description, completed, due_date, repeat_type, created_at, updated_at FROM todos WHERE id = ?`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, svcerror.ErrTodoNotFound
		}
		return nil, err
	}
	return &t, nil
}

// List query fragments — all fully parameterized; sort selects which query to run.
const (
	_listCount = `SELECT COUNT(*) FROM todos WHERE (? = 1 OR completed = 0)`
	_listBase  = `SELECT id, title, description, completed, due_date, repeat_type, created_at, updated_at FROM todos WHERE (? = 1 OR completed = 0)`
)

var listQueries = map[string]string{
	"due_date_desc":   _listBase + ` ORDER BY due_date IS NULL DESC, due_date DESC, id DESC LIMIT ? OFFSET ?`,
	"created_at_asc": _listBase + ` ORDER BY created_at ASC, id ASC LIMIT ? OFFSET ?`,
	"created_at_desc": _listBase + ` ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`,
	"due_date_asc":   _listBase + ` ORDER BY due_date IS NULL, due_date ASC, id ASC LIMIT ? OFFSET ?`,
}

func (r *todoRepository) List(ctx context.Context, includeCompleted bool, sort string, limit, offset int) ([]*todo.Todo, int, error) {
	includeAll := 0
	if includeCompleted {
		includeAll = 1
	}
	query, ok := listQueries[sort]
	if !ok {
		query = listQueries["due_date_asc"]
	}

	var total int
	err := r.db.GetContext(ctx, &total, _listCount, includeAll)
	if err != nil {
		return nil, 0, err
	}

	var items []*todo.Todo
	err = r.db.SelectContext(ctx, &items, query, includeAll, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *todoRepository) Update(ctx context.Context, id int64, title, description *string, completed *bool, dueDate *time.Time, repeatType *string) (*todo.Todo, error) {
	var set []string
	var args []interface{}
	if title != nil {
		set = append(set, "title = ?")
		args = append(args, *title)
	}
	if description != nil {
		set = append(set, "description = ?")
		args = append(args, *description)
	}
	if completed != nil {
		set = append(set, "completed = ?")
		args = append(args, *completed)
	}
	if dueDate != nil {
		set = append(set, "due_date = ?")
		args = append(args, *dueDate)
	}
	if repeatType != nil {
		set = append(set, "repeat_type = ?")
		args = append(args, *repeatType)
	}
	if len(set) == 0 {
		return r.GetByID(ctx, id)
	}
	set = append(set, "updated_at = CURRENT_TIMESTAMP(3)")
	args = append(args, id)
	query := fmt.Sprintf(`UPDATE todos SET %s WHERE id = ?`, strings.Join(set, ", "))
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, svcerror.ErrTodoNotFound
	}
	return r.GetByID(ctx, id)
}

func (r *todoRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM todos WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return svcerror.ErrTodoNotFound
	}
	return nil
}
