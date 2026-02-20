package upgrade

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sanketmote/gokit-wrapper/logger/svclog"
)

// Stages is the ordered list of incremental upgrade stages.
var Stages = []Stage{
	{Version: Version, Fn: Upgrade1},
}

// Run executes DB upgrade: Install for fresh DB, Incremental for existing.
func Run(ctx context.Context, db *sqlx.DB, logger svclog.Logger) error {
	fresh, err := IsFreshDB(ctx, db)
	if err != nil {
		return err
	}

	u := NewUpgrader(db, logger)

	if fresh {
		logger.Info(ctx, "type", "install", "msg", "running install for fresh DB")
		if err := u.Install(ctx, InstallAll); err != nil {
			return err
		}
		logger.Info(ctx, "type", "install", "msg", "install completed")
		return nil
	}

	if err := CheckUpgradeInProgress(ctx, db); err != nil {
		return err
	}

	logger.Info(ctx, "type", "incremental", "msg", "running incremental upgrade")
	if err := u.Incremental(ctx, Stages); err != nil {
		return err
	}
	logger.Info(ctx, "type", "incremental", "msg", "incremental upgrade completed")
	return nil
}
