package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sanketmote/go-todo-service/internal/svcerror"
)

func newTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "mysql"), mock
}

func TestTodoRepository_Create(t *testing.T) {
	ctx := context.Background()
	due := time.Date(2025, 2, 25, 18, 0, 0, 0, time.UTC)

	t.Run("ok", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		mock.ExpectExec(`INSERT INTO todos \(title, description, completed, due_date, repeat_type\) VALUES \(\?, \?, 0, \?, \?\)`).
			WithArgs("Buy groceries", "Milk", due, "none").
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"id", "title", "description", "completed", "due_date", "repeat_type", "created_at", "updated_at"}).
			AddRow(1, "Buy groceries", "Milk", false, due, "none", time.Now(), time.Now())
		mock.ExpectQuery(`SELECT id, title, description, completed, due_date, repeat_type, created_at, updated_at FROM todos WHERE id = \?`).
			WithArgs(1).
			WillReturnRows(rows)

		repo := NewTodoRepository(sqlxDB)
		got, err := repo.Create(ctx, "Buy groceries", "Milk", &due, "none")
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if got.ID != 1 || got.Title != "Buy groceries" {
			t.Errorf("got %+v", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("expectations: %v", err)
		}
	})

	t.Run("empty repeat defaults to none", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		mock.ExpectExec(`INSERT INTO todos`).
			WithArgs("Task", "", nil, "none").
			WillReturnResult(sqlmock.NewResult(1, 1))

		rows := sqlmock.NewRows([]string{"id", "title", "description", "completed", "due_date", "repeat_type", "created_at", "updated_at"}).
			AddRow(1, "Task", "", false, nil, "none", time.Now(), time.Now())
		mock.ExpectQuery(`SELECT`).
			WithArgs(1).
			WillReturnRows(rows)

		repo := NewTodoRepository(sqlxDB)
		got, err := repo.Create(ctx, "Task", "", nil, "")
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if got.RepeatType != "none" {
			t.Errorf("RepeatType = %q, want none", got.RepeatType)
		}
	})
}

func TestTodoRepository_GetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		rows := sqlmock.NewRows([]string{"id", "title", "description", "completed", "due_date", "repeat_type", "created_at", "updated_at"}).
			AddRow(1, "Task", "", false, nil, "none", time.Now(), time.Now())
		mock.ExpectQuery(`SELECT id, title, description, completed, due_date, repeat_type, created_at, updated_at FROM todos WHERE id = \?`).
			WithArgs(1).
			WillReturnRows(rows)

		repo := NewTodoRepository(sqlxDB)
		got, err := repo.GetByID(ctx, 1)
		if err != nil {
			t.Fatalf("GetByID: %v", err)
		}
		if got.ID != 1 {
			t.Errorf("ID = %d, want 1", got.ID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		mock.ExpectQuery(`SELECT`).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		repo := NewTodoRepository(sqlxDB)
		_, err := repo.GetByID(ctx, 999)
		if !errors.Is(err, svcerror.ErrTodoNotFound) {
			t.Errorf("err = %v, want ErrTodoNotFound", err)
		}
	})
}

func TestTodoRepository_List(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		mock.ExpectQuery(`SELECT COUNT\(\*\) FROM todos`).
			WithArgs(0).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		rows := sqlmock.NewRows([]string{"id", "title", "description", "completed", "due_date", "repeat_type", "created_at", "updated_at"}).
			AddRow(1, "A", "", false, nil, "none", time.Now(), time.Now()).
			AddRow(2, "B", "", false, nil, "none", time.Now(), time.Now())
		mock.ExpectQuery(`SELECT id, title, description, completed, due_date, repeat_type, created_at, updated_at FROM todos`).
			WithArgs(0, 50, 0).
			WillReturnRows(rows)

		repo := NewTodoRepository(sqlxDB)
		items, total, err := repo.List(ctx, false, "due_date_asc", 50, 0)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if total != 2 || len(items) != 2 {
			t.Errorf("total=%d len=%d, want 2,2", total, len(items))
		}
	})
}

func TestTodoRepository_Update(t *testing.T) {
	ctx := context.Background()
	title := "Updated"

	t.Run("ok", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		mock.ExpectExec(`UPDATE todos SET`).
			WithArgs("Updated", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		rows := sqlmock.NewRows([]string{"id", "title", "description", "completed", "due_date", "repeat_type", "created_at", "updated_at"}).
			AddRow(1, "Updated", "", false, nil, "none", time.Now(), time.Now())
		mock.ExpectQuery(`SELECT`).
			WithArgs(1).
			WillReturnRows(rows)

		repo := NewTodoRepository(sqlxDB)
		got, err := repo.Update(ctx, 1, &title, nil, nil, nil, nil)
		if err != nil {
			t.Fatalf("Update: %v", err)
		}
		if got.Title != "Updated" {
			t.Errorf("Title = %q, want Updated", got.Title)
		}
	})

	t.Run("not found", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		mock.ExpectExec(`UPDATE todos SET`).
			WithArgs("Updated", 999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewTodoRepository(sqlxDB)
		_, err := repo.Update(ctx, 999, &title, nil, nil, nil, nil)
		if !errors.Is(err, svcerror.ErrTodoNotFound) {
			t.Errorf("err = %v, want ErrTodoNotFound", err)
		}
	})
}

func TestTodoRepository_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		mock.ExpectExec(`DELETE FROM todos WHERE id = \?`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo := NewTodoRepository(sqlxDB)
		err := repo.Delete(ctx, 1)
		if err != nil {
			t.Fatalf("Delete: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		sqlxDB, mock := newTestDB(t)
		defer sqlxDB.Close()

		mock.ExpectExec(`DELETE FROM todos WHERE id = \?`).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		repo := NewTodoRepository(sqlxDB)
		err := repo.Delete(ctx, 999)
		if !errors.Is(err, svcerror.ErrTodoNotFound) {
			t.Errorf("err = %v, want ErrTodoNotFound", err)
		}
	})
}

