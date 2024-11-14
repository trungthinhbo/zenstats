package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/dreamsofcode-io/zenstats/internal/quote"
	"github.com/dreamsofcode-io/zenstats/internal/repository"
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
	repo   *repository.Queries
}

func NewIndex(
	logger *slog.Logger,
	quotes *quote.Service,
	tmpl *template.Template,
	repo *repository.Queries,
) *Index {
	index := &Index{
		logger: logger,
		tmpl:   tmpl,
		quotes: quotes,
		repo:   repo,
	}

	return index
}

type Visitors struct {
	LastHour int64
	LastDay  int64
	LastWeek int64
}

type pageData struct {
	Quote     quote.Quote
	TotalHits int64
	Visitors  Visitors
}

var hits = 1057

var re = regexp.MustCompile(`,\s*`)

// TODO - Test this
func getIP(r *http.Request, steps int) string {
	xff := r.Header.Get("X-Forwarded-For")
	xffParts := re.Split(xff, -1)

	if xff == "" || len(xffParts) <= steps {
		remoteParts := strings.Split(r.RemoteAddr, ":")
		return strings.Join(remoteParts[:len(remoteParts)-1], ":")
	}

	return xffParts[len(xffParts)-(steps+1)]
}

func (h *Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r, 0)

	ctx := r.Context()

	if err := h.repo.InsertVisit(ctx, ip); err != nil {
		h.logger.Error("failed to insert visit", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	total, err := h.repo.CountAllVisits(ctx)
	if err != nil {
		h.logger.Error("failed to get total visits", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hour, err := h.repo.CountVisitors(ctx, time.Now().Add(-time.Hour))
	if err != nil {
		h.logger.Error("failed to get hour visits", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	day, err := h.repo.CountVisitors(ctx, time.Now().Add(-time.Hour*24))
	if err != nil {
		h.logger.Error("failed to get day visits", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	week, err := h.repo.CountVisitors(ctx, time.Now().Add(-time.Hour*24*7))
	if err != nil {
		h.logger.Error("failed to get week visits", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	quote := h.quotes.GetQuote()
	h.tmpl.ExecuteTemplate(w, "index.html", pageData{
		Quote:     quote,
		TotalHits: total,
		Visitors: Visitors{
			LastHour: hour,
			LastDay:  day,
			LastWeek: week,
		},
	})
}
