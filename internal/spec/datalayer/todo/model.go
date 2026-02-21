package todo

import "time"

// Todo represents a row in the todos table (DB mapping).
type Todo struct {
	ID          int64      `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Completed   bool       `db:"completed"`
	DueDate     *time.Time `db:"due_date"`
	RepeatType  string     `db:"repeat_type"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
