package handler

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/dreamsofcode-io/zenfulstats/internal/quote"
)

func Page(tmpl *template.Template, name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, name, nil)
	}
}

type Index struct {
	logger *slog.Logger
	tmpl   *template.Template
	quotes *quote.Service
}

func NewIndex(logger *slog.Logger, quotes *quote.Service, tmpl *template.Template) *Index {
	index := &Index{
		logger: logger,
		tmpl:   tmpl,
		quotes: quotes,
	}

	return index
}

type Visitors struct {
	LastHour int
	LastDay  int
	LastWeek int
}

type pageData struct {
	Quote     quote.Quote
	TotalHits int
	Visitors  Visitors
}

var hits = 1057

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	quote := h.quotes.GetQuote()
	hits += 1
	h.tmpl.ExecuteTemplate(w, "index.html", pageData{
		Quote:     quote,
		TotalHits: hits,
		Visitors: Visitors{
			LastHour: 27,
			LastDay:  312,
			LastWeek: 2311,
		},
	})
}
