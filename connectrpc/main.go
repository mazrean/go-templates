package main

import (
	"log/slog"
	"os"

	"github.com/mazrean/go-templates/connectrpc/internal/di"
)

func main() {
	app, err := di.DI()
	if err != nil {
		slog.Error(
			"failed to inject app",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		slog.Error(
			"failed to run app",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
