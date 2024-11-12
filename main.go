//go:generate tailwindcss -i ./static/css/input.css -o ./static/css/style.css --minify
package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"os/signal"

	"github.com/dreamsofcode-io/zenstats/internal/app"
)

//go:embed templates/*.html static
var files embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := app.New(logger, app.Config{}, files)

	if err := app.Start(ctx); err != nil {
		logger.Error("failed to start app", slog.Any("error", err))
	}
}
