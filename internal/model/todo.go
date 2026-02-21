package model

import "time"

// RepeatType defines how often a todo recurs. Empty or "none" means not recurring.
const (
	RepeatNone    = "none"
	RepeatDaily   = "daily"
	RepeatWeekly  = "weekly"
	RepeatMonthly = "monthly"
	RepeatYearly  = "yearly"
)

// Todo is the domain model for a todo item (DB and in-memory).
// API uses snake_case for JSON fields per LLD.
type Todo struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"` // Optional longer text; empty means none
	Completed   bool       `db:"completed"`
	DueDate     *time.Time `db:"due_date"`
	RepeatType  string     `db:"repeat_type"` // none, daily, weekly, monthly, yearly
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// CreateTodoRequest is the request body for POST /api/v1/todos.
type CreateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"` // Optional
	DueDate     *time.Time `json:"due_date,omitempty"`
	RepeatType  string     `json:"repeat_type,omitempty"` // none, daily, weekly, monthly, yearly
}

// UpdateTodoRequest is the request body for PUT /api/v1/todos/:id (partial update).
type UpdateTodoRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Completed   *bool      `json:"completed,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	RepeatType  *string    `json:"repeat_type,omitempty"`
}

// TodoResponse is the JSON response for a single todo.
type TodoResponse struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Completed   bool       `json:"completed"`
	DueDate     *time.Time `json:"due_date"`              // null when no due date
	RepeatType  string     `json:"repeat_type,omitempty"` // none, daily, weekly, monthly, yearly
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ListTodosResponse is the JSON response for GET /api/v1/todos.
type ListTodosResponse struct {
	Items []TodoResponse `json:"items"`
	Total int64          `json:"total"`
}

// ToTodoResponse converts a domain Todo to TodoResponse.
func ToTodoResponse(t *Todo) TodoResponse {
	return TodoResponse{
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
