package todo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sanketmote/go-todo-service/internal/repository"
	dltodo "github.com/sanketmote/go-todo-service/internal/spec/datalayer/todo"
	spectodo "github.com/sanketmote/go-todo-service/internal/spec/todo"
	"github.com/sanketmote/go-todo-service/internal/svcerror"
)

// Ensure mockRepo implements repository.TodoRepository.
var _ repository.TodoRepository = (*mockRepo)(nil)

// mockRepo implements repository.TodoRepository for tests.
type mockRepo struct {
	createFunc  func(ctx context.Context, title, desc string, due *time.Time, repeat string) (*dltodo.Todo, error)
	getByIDFunc func(ctx context.Context, id int64) (*dltodo.Todo, error)
	listFunc    func(ctx context.Context, inc bool, sort string, limit, offset int) ([]*dltodo.Todo, int, error)
	updateFunc  func(ctx context.Context, id int64, title, desc *string, completed *bool, due *time.Time, repeat *string) (*dltodo.Todo, error)
	deleteFunc  func(ctx context.Context, id int64) error
}

func (m *mockRepo) Create(ctx context.Context, title, description string, dueDate *time.Time, repeatType string) (*dltodo.Todo, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, title, description, dueDate, repeatType)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*dltodo.Todo, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) List(ctx context.Context, includeCompleted bool, sort string, limit, offset int) ([]*dltodo.Todo, int, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, includeCompleted, sort, limit, offset)
	}
	return nil, 0, errors.New("not implemented")
}

func (m *mockRepo) Update(ctx context.Context, id int64, title, description *string, completed *bool, dueDate *time.Time, repeatType *string) (*dltodo.Todo, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, title, description, completed, dueDate, repeatType)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return errors.New("not implemented")
}

func TestTodoService_Create(t *testing.T) {
	ctx := context.Background()
	due := time.Date(2025, 2, 25, 18, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		title   string
		desc    string
		due     *time.Time
		repeat  string
		repo    *mockRepo
		wantErr error
	}{
		{
			name:   "ok",
			title:  "Buy groceries",
			desc:   "Milk, eggs",
			due:    &due,
			repeat: "none",
			repo: &mockRepo{
				createFunc: func(_ context.Context, title, desc string, due *time.Time, repeat string) (*dltodo.Todo, error) {
					return &dltodo.Todo{ID: 1, Title: title, Description: desc, DueDate: due, RepeatType: repeat, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				},
			},
			wantErr: nil,
		},
		{
			name:    "empty title",
			title:   "",
			desc:    "",
			due:     nil,
			repeat:  "none",
			repo:    &mockRepo{},
			wantErr: &svcerror.ValidationError{},
		},
		{
			name:    "whitespace title",
			title:   "   ",
			desc:    "",
			due:     nil,
			repeat:  "none",
			repo:    &mockRepo{},
			wantErr: &svcerror.ValidationError{},
		},
		{
			name:    "title too long",
			title:   string(make([]byte, 501)),
			desc:    "",
			due:     nil,
			repeat:  "none",
			repo:    &mockRepo{},
			wantErr: &svcerror.ValidationError{},
		},
		{
			name:    "invalid repeat_type",
			title:   "Task",
			desc:    "",
			due:     nil,
			repeat:  "invalid",
			repo:    &mockRepo{},
			wantErr: &svcerror.ValidationError{},
		},
		{
			name:   "empty repeat accepted",
			title:  "Task",
			desc:   "",
			due:    nil,
			repeat: "",
			repo: &mockRepo{
				createFunc: func(_ context.Context, title, _ string, _ *time.Time, repeat string) (*dltodo.Todo, error) {
					rt := repeat
					if rt == "" {
						rt = "none"
					}
					return &dltodo.Todo{ID: 1, Title: title, RepeatType: rt, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTodoService(tt.repo, nil)
			got, err := svc.Create(ctx, tt.title, tt.desc, tt.due, tt.repeat)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var vErr *svcerror.ValidationError
				if !errors.As(err, &vErr) {
					t.Errorf("err = %v, want ValidationError", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.Title != tt.title {
				t.Errorf("Title = %q, want %q", got.Title, tt.title)
			}
		})
	}
}

func TestTodoService_GetByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		id      int64
		repo    *mockRepo
		wantErr error
	}{
		{
			name: "ok",
			id:   1,
			repo: &mockRepo{
				getByIDFunc: func(_ context.Context, id int64) (*dltodo.Todo, error) {
					return &dltodo.Todo{ID: id, Title: "Task", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				},
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   999,
			repo: &mockRepo{
				getByIDFunc: func(_ context.Context, _ int64) (*dltodo.Todo, error) {
					return nil, svcerror.ErrTodoNotFound
				},
			},
			wantErr: svcerror.ErrTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTodoService(tt.repo, nil)
			got, err := svc.GetByID(ctx, tt.id)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.ID != tt.id {
				t.Errorf("ID = %d, want %d", got.ID, tt.id)
			}
		})
	}
}

func TestTodoService_List(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		include    bool
		sort       string
		limit      int
		offset     int
		repo       *mockRepo
		wantErr    bool
		wantTotal  int
		wantItems  int
	}{
		{
			name:      "ok",
			include:   false,
			sort:      spectodo.SortDueDateAsc,
			limit:     10,
			offset:    0,
			wantTotal:  2,
			wantItems:  2,
			repo: &mockRepo{
				listFunc: func(_ context.Context, _ bool, _ string, limit, offset int) ([]*dltodo.Todo, int, error) {
					items := []*dltodo.Todo{
						{ID: 1, Title: "A", CreatedAt: time.Now(), UpdatedAt: time.Now()},
						{ID: 2, Title: "B", CreatedAt: time.Now(), UpdatedAt: time.Now()},
					}
					return items, 2, nil
				},
			},
		},
		{
			name:    "invalid sort",
			include: false,
			sort:    "invalid_sort",
			limit:   10,
			offset:  0,
			repo:    &mockRepo{},
			wantErr: true,
		},
		{
			name:      "empty sort defaults to due_date_asc",
			include:   false,
			sort:      "",
			limit:     50,
			offset:    0,
			wantTotal: 0,
			wantItems: 0,
			repo: &mockRepo{
				listFunc: func(_ context.Context, _ bool, sort string, limit, offset int) ([]*dltodo.Todo, int, error) {
					if sort != spectodo.SortDueDateAsc {
						return nil, 0, errors.New("expected due_date_asc")
					}
					return nil, 0, nil
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTodoService(tt.repo, nil)
			got, err := svc.List(ctx, tt.include, tt.sort, tt.limit, tt.offset)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got.Total != int64(tt.wantTotal) {
				t.Errorf("Total = %d, want %d", got.Total, tt.wantTotal)
			}
			if len(got.Items) != tt.wantItems {
				t.Errorf("len(Items) = %d, want %d", len(got.Items), tt.wantItems)
			}
		})
	}
}

func TestTodoService_Update(t *testing.T) {
	ctx := context.Background()
	title := "Updated"
	completed := true

	tests := []struct {
		name    string
		id      int64
		title   *string
		desc    *string
		completed *bool
		repo    *mockRepo
		wantErr error
	}{
		{
			name:      "ok",
			id:        1,
			title:     &title,
			completed: &completed,
			repo: &mockRepo{
				updateFunc: func(_ context.Context, id int64, title, _ *string, completed *bool, _ *time.Time, _ *string) (*dltodo.Todo, error) {
					return &dltodo.Todo{ID: id, Title: *title, Completed: *completed, CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
				},
			},
			wantErr: nil,
		},
		{
			name:    "empty title",
			id:      1,
			title:   strPtr("   "),
			repo:    &mockRepo{},
			wantErr: &svcerror.ValidationError{},
		},
		{
			name:    "not found",
			id:      999,
			title:   &title,
			repo: &mockRepo{
				updateFunc: func(_ context.Context, _ int64, _ *string, _ *string, _ *bool, _ *time.Time, _ *string) (*dltodo.Todo, error) {
					return nil, svcerror.ErrTodoNotFound
				},
			},
			wantErr: svcerror.ErrTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTodoService(tt.repo, nil)
			got, err := svc.Update(ctx, tt.id, tt.title, tt.desc, tt.completed, nil, nil)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				var vErr *svcerror.ValidationError
				if !errors.Is(err, tt.wantErr) && !errors.As(err, &vErr) {
					t.Errorf("err = %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if tt.title != nil && got.Title != *tt.title {
				t.Errorf("Title = %q, want %q", got.Title, *tt.title)
			}
		})
	}
}

func TestTodoService_Delete(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		id      int64
		repo    *mockRepo
		wantErr error
	}{
		{
			name: "ok",
			id:   1,
			repo: &mockRepo{
				deleteFunc: func(_ context.Context, _ int64) error { return nil },
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   999,
			repo: &mockRepo{
				deleteFunc: func(_ context.Context, _ int64) error { return svcerror.ErrTodoNotFound },
			},
			wantErr: svcerror.ErrTodoNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewTodoService(tt.repo, nil)
			err := svc.Delete(ctx, tt.id)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
