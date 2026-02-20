package dbversion

// DBVersion represents the single row in db_version (id=1).
type DBVersion struct {
	ID        int64 `json:"id" db:"id"`
	Version   int64 `json:"version" db:"version"`
	Inupgrade bool  `json:"inupgrade" db:"inupgrade"`
}

// DBVersionHistory represents a row in db_version_history.
type DBVersionHistory struct {
	ID          int64  `json:"id" db:"id"`
	Version     int64  `json:"version" db:"version"`
	Timestamp   int64  `json:"timestamp" db:"timestamp"`
	Description string `json:"description" db:"description"`
}
