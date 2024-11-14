package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/dreamsofcode-io/zenstats/internal/database"
	"github.com/dreamsofcode-io/zenstats/internal/quote"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App contains all of the application dependencies for the project.
type App struct {
	config Config
	logger *slog.Logger
	files  fs.FS
	quotes *quote.Service
	db     *pgxpool.Pool
}

// New creates a new instance of the application.
func New(logger *slog.Logger, config Config, files fs.FS) *App {
	return &App{
		config: config,
		logger: logger,
		files:  files,
		quotes: quote.New(),
	}
}

// Start is used to start the application. The application
// will run until either the given context is cancelled, or
// the application is ended.
func (a *App) Start(ctx context.Context) error {
	db, err := database.Connect(ctx, a.logger, a.files)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	a.db = db

	router, err := a.loadRoutes()
	if err != nil {
		return fmt.Errorf("failed when loading routes: %w", err)
	}

	port := 8080
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	errCh := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to listen and serve: %w", err)
		}

		close(errCh)
	}()

	a.logger.Info("server running", slog.Int("port", port))

	select {
	// Wait until we receive SIGINT (ctrl+c on cli)
	case <-ctx.Done():
		break
	case err := <-errCh:
		return err
	}

	sCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(sCtx)

	return nil
}
