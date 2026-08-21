package main

import (
	"sync"

	"github.com/z46-dev/florr-legacy/command/app"
	"github.com/z46-dev/florr-legacy/command/config"
	"github.com/z46-dev/florr-legacy/command/db"
	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/golog"
)

func main() {
	var err error

	shared.L = golog.New().Prefix("[COMMAND]", golog.BoldCyan).Timestamp()
	if err = config.Init("config.toml"); err != nil {
		shared.L.Fatalf("Failed to initialize config: %v\n", err)
		return
	}

	if err = db.InitDB(); err != nil {
		shared.L.Fatalf("Failed to initialize database: %v\n", err)
		return
	}

	defer func() {
		if err = db.Close(); err != nil {
			shared.L.Errorf("Failed to close database: %v\n", err)
		}
	}()

	var wg sync.WaitGroup

	wg.Go(func() {
		if err = app.Start(); err != nil {
			shared.L.Fatalf("Web server stopped unexpectedly: %v\n", err)
		}
	})

	shared.L.Info("Command service started successfully")

	wg.Wait()
	shared.L.Info("Command service stopped")
}
