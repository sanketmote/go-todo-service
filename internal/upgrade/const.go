package upgrade

const (
	// VersionTableName is the name of the schema version table.
	VersionTableName = "db_version"

	// VersionHistoryTableName is the name of the version history table.
	VersionHistoryTableName = "db_version_history"

	// Version is the current schema version (int64).
	Version = 1

	// InsertVersionQuery inserts the single version row (id=1). Args: version.
	InsertVersionQuery = `INSERT INTO db_version (id, version, inupgrade) VALUES (1, ?, 0)`

	// SetInupgradeQuery sets inupgrade=1 before starting upgrade.
	SetInupgradeQuery = `UPDATE db_version SET inupgrade = 1 WHERE id = 1`

	// ClearInupgradeQuery sets inupgrade=0 and updates version after upgrade. Args: version.
	ClearInupgradeQuery = `UPDATE db_version SET version = ?, inupgrade = 0 WHERE id = 1`

	// GetVersionQuery returns version and inupgrade for the single row.
	GetVersionQuery = `SELECT version, inupgrade FROM db_version WHERE id = 1 LIMIT 1`

	// InsertVersionHistoryQuery records an upgrade in history. Args: version, timestamp, description.
	InsertVersionHistoryQuery = `INSERT INTO db_version_history (version, timestamp, description) VALUES (?, ?, ?)`
)
