package db

import (
	"os"
	"slices"

	"github.com/z46-dev/florr-legacy/command/config"
	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/gosqlite"
)

var (
	driver           *gosqlite.Driver
	Accounts         *gosqlite.RegisteredStruct[shared.Account]
	AuthMethods      *gosqlite.RegisteredStruct[shared.AuthMethod]
	Discords         *gosqlite.RegisteredStruct[shared.Discord]
	DiscordHistories *gosqlite.RegisteredStruct[shared.DiscordHistory]
)

func migrateTable[T any](table *gosqlite.RegisteredStruct[T], migrationOpts gosqlite.MigrationOptions) (err error) {
	var report *gosqlite.MigrationReport
	if report, err = table.Migrate(migrationOpts); err == nil && report != nil {
		if len(report.AddedColumns) > 0 {
			shared.L.Warningf("Added columns to table for %T: %v\n", shared.Account{}, report.AddedColumns)
		}

		if len(report.ChangedColumns) > 0 {
			shared.L.Warningf("Changed columns in table for %T: %v\n", shared.Account{}, report.ChangedColumns)
		}

		if len(report.DroppedColumns) > 0 {
			shared.L.Warningf("Dropped columns from table for %T: %v\n", shared.Account{}, report.DroppedColumns)
		}

		if len(report.RenamedColumns) > 0 {
			shared.L.Warningf("Renamed columns in table for %T: %v\n", shared.Account{}, report.RenamedColumns)
		}

		if report.Rebuilt {
			shared.L.Warningf("Rebuilt table for %T\n", shared.Account{})
		}
	}

	return
}

func InitDB() (err error) {
	if driver, err = gosqlite.Begin(config.Config.Database.File); err != nil {
		return
	}

	if Accounts, err = gosqlite.Register(driver, shared.Account{}); err != nil {
		return
	}

	if AuthMethods, err = gosqlite.Register(driver, shared.AuthMethod{}); err != nil {
		return
	}

	if Discords, err = gosqlite.Register(driver, shared.Discord{}); err != nil {
		return
	}

	if DiscordHistories, err = gosqlite.Register(driver, shared.DiscordHistory{}); err != nil {
		return
	}

	var migrationOpts gosqlite.MigrationOptions

	if len(os.Args) > 1 && slices.Contains(os.Args, "--allow-destructive-migrations") {
		migrationOpts.AllowDestructive = true
		shared.L.Warning("Destructive migrations are enabled!")
	}

	if err = migrateTable(Accounts, migrationOpts); err != nil {
		return
	}

	if err = migrateTable(AuthMethods, migrationOpts); err != nil {
		return
	}

	if err = migrateTable(Discords, migrationOpts); err != nil {
		return
	}

	if err = migrateTable(DiscordHistories, migrationOpts); err != nil {
		return
	}

	return
}

func Close() (err error) {
	if driver != nil {
		err = driver.Close()
	}

	return
}
