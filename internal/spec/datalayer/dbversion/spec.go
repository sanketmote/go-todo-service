package dbversion

const (
	// CreateDBVersionTableQuery creates a single-row table for schema version.
	// Application must maintain exactly one row (id=1).
	CreateDBVersionTableQuery = `
CREATE TABLE IF NOT EXISTS db_version (
	id BIGINT PRIMARY KEY DEFAULT 1,
	version BIGINT NOT NULL DEFAULT 0,
	inupgrade TINYINT(1) NOT NULL DEFAULT 0,
	CONSTRAINT chk_single_row CHECK (id = 1)
);
`

	CreateDBVersionHistoryTableQuery = `
CREATE TABLE IF NOT EXISTS db_version_history (
	id BIGINT AUTO_INCREMENT PRIMARY KEY,
	version BIGINT NOT NULL,
	timestamp BIGINT NOT NULL,
	description VARCHAR(500) NOT NULL DEFAULT '',
	INDEX idx_version (version)
);
`
)
