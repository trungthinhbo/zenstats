package app

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/dreamsofcode-io/zenstats/internal/handler"
	"github.com/dreamsofcode-io/zenstats/internal/middleware"
	"github.com/dreamsofcode-io/zenstats/internal/repository"
)

const index = "GET /{$}"

// TODO - This function is where your pages are loaded.
func (a *App) loadPages(router *http.ServeMux, tmpl *template.Template) {
	repo := repository.New(a.db)

	h := handler.NewIndex(a.logger, a.quotes, tmpl, repo)

	// This is your index route, i.e. /. It has an odd syntax in
	// the go serve mux
	router.Handle("GET /{$}", h)
}

func (a *App) loadRoutes() (http.Handler, error) {
	tmpl, err := template.New("").ParseFS(a.files, "templates/*")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	static, err := fs.Sub(a.files, "static")
	if err != nil {
		return nil, fmt.Errorf("failed to subdir static: %w", err)
	}

	// Create a new router
	router := http.NewServeMux()

	// This is the static fileserver.
	router.Handle("GET /static/", http.StripPrefix("/static", http.FileServerFS(static)))

	a.loadPages(router, tmpl)

	// Create a middleware chain from the Chain function of the
	// middleware package
	chain := middleware.Chain(
		middleware.Logging(a.logger),
	)

	return chain(router), nil
}
