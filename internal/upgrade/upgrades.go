package upgrade

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sanketmote/go-todo-service/internal/spec/datalayer/dbversion"
	"github.com/sanketmote/go-todo-service/internal/spec/datalayer/todo"
)

// InstallAll creates the version table, version history table, and todos table for a fresh DB.
func InstallAll(ctx context.Context, tx *sqlx.Tx) error {
	if _, err := tx.ExecContext(ctx, dbversion.CreateDBVersionTableQuery); err != nil {
		return fmt.Errorf("create db_version: %w", err)
	}
	if _, err := tx.ExecContext(ctx, dbversion.CreateDBVersionHistoryTableQuery); err != nil {
		return fmt.Errorf("create db_version_history: %w", err)
	}
	if _, err := tx.ExecContext(ctx, todo.CreateTodosTableQuery); err != nil {
		return fmt.Errorf("create todos: %w", err)
	}
	if _, err := tx.ExecContext(ctx, InsertVersionQuery, Version); err != nil {
		return fmt.Errorf("insert version: %w", err)
	}
	return nil
}

// Upgrade1 brings a DB with version table but no todos to version 1.
// Creates todos table; version row is updated by Upgrader.Incremental.
func Upgrade1(ctx context.Context, tx *sqlx.Tx) error {
	_, err := tx.ExecContext(ctx, todo.CreateTodosTableQuery)
	return err
}
