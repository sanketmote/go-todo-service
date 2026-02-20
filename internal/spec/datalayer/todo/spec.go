package todo

const (
	CreateTodosTableQuery = `
CREATE TABLE IF NOT EXISTS todos (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	title VARCHAR(500) NOT NULL,
	description TEXT,
	completed TINYINT(1) NOT NULL DEFAULT 0,
	due_date DATETIME(3) NULL,
	repeat_type VARCHAR(32) NOT NULL DEFAULT 'none',
	created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
	updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
	INDEX idx_list (completed, due_date, id)
);
`
)
