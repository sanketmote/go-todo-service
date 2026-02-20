package upgrade

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sanketmote/gokit-wrapper/logger/svclog"
)

// ErrUpgradeInProgress is returned when a previous upgrade did not complete (inupgrade=1).
var ErrUpgradeInProgress = errors.New("upgrade already in progress; previous upgrade may have crashed")

// Upgrader runs DB schema upgrades (Install or Incremental).
type Upgrader struct {
	db     *sqlx.DB
	logger svclog.Logger
}

// NewUpgrader returns a new Upgrader.
func NewUpgrader(db *sqlx.DB, logger svclog.Logger) *Upgrader {
	return &Upgrader{db: db, logger: logger}
}

// IsFreshDB returns true if the version table does not exist or has no row.
func IsFreshDB(ctx context.Context, db *sqlx.DB) (bool, error) {
	var exists int
	err := db.GetContext(ctx, &exists,
		"SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?",
		VersionTableName)
	if err != nil {
		return false, err
	}
	if exists == 0 {
		return true, nil
	}

	var count int
	err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM "+VersionTableName)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// CheckUpgradeInProgress returns ErrUpgradeInProgress if inupgrade=1.
func CheckUpgradeInProgress(ctx context.Context, db *sqlx.DB) error {
	var version int64
	var inupgrade int
	err := db.QueryRowContext(ctx, GetVersionQuery).Scan(&version, &inupgrade)
	if err != nil {
		return err
	}
	if inupgrade == 1 {
		return ErrUpgradeInProgress
	}
	return nil
}

// Install creates the version table and all app tables for a fresh DB.
func (u *Upgrader) Install(ctx context.Context, installFn func(context.Context, *sqlx.Tx) error) error {
	tx, err := u.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := installFn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}

// Incremental runs upgrade stages from current version to latest.
// Sets inupgrade=1 before starting, clears it after each stage.
func (u *Upgrader) Incremental(ctx context.Context, stages []Stage) error {
	var current int64
	var inupgrade int
	err := u.db.QueryRowContext(ctx, GetVersionQuery).Scan(&current, &inupgrade)
	if err != nil {
		return fmt.Errorf("get current version: %w", err)
	}
	if inupgrade == 1 {
		return ErrUpgradeInProgress
	}

	for _, s := range stages {
		if s.Version <= current {
			continue
		}

		// Set inupgrade=1 before starting
		_, err = u.db.ExecContext(ctx, SetInupgradeQuery)
		if err != nil {
			return fmt.Errorf("set inupgrade: %w", err)
		}

		tx, err := u.db.BeginTxx(ctx, nil)
		if err != nil {
			_, _ = u.db.ExecContext(ctx, ClearInupgradeQuery, current)
			return err
		}

		if err := s.Fn(ctx, tx); err != nil {
			_ = tx.Rollback()
			_, _ = u.db.ExecContext(ctx, ClearInupgradeQuery, current)
			return err
		}

		if err := tx.Commit(); err != nil {
			_, _ = u.db.ExecContext(ctx, ClearInupgradeQuery, current)
			return err
		}

		// Clear inupgrade and update version (DDL in stage may have committed; do this outside stage tx)
		_, err = u.db.ExecContext(ctx, ClearInupgradeQuery, s.Version)
		if err != nil {
			return fmt.Errorf("clear inupgrade: %w", err)
		}

		_, err = u.db.ExecContext(ctx, InsertVersionHistoryQuery, s.Version, time.Now().Unix(), fmt.Sprintf("upgrade to v%d", s.Version))
		if err != nil {
			return fmt.Errorf("insert version history: %w", err)
		}

		current = s.Version
	}
	return nil
}

// Stage is a single upgrade step (version + function).
type Stage struct {
	Version int64
	Fn      func(context.Context, *sqlx.Tx) error
}
