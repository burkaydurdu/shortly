package server

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	shortlyLog "github.com/burkaydurdu/shortly/pkg/log"

	shortlyDB "github.com/burkaydurdu/shortly/internal/db"

	"github.com/burkaydurdu/shortly/internal/domain/shortly"

	"github.com/burkaydurdu/shortly/config"
	shortlyError "github.com/burkaydurdu/shortly/pkg/error"
	httpSwagger "github.com/swaggo/http-swagger"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type ShortlyMux struct {
	routes []*route
	l      *shortlyLog.ShortlyLog
}

func (h *ShortlyMux) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *ShortlyMux) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *ShortlyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}

	h.l.Zap(fmt.Sprintf("%s %s %s", r.Method, r.URL, shortlyError.PathNotFoundError))
	http.NotFound(w, r)
}

type Server struct {
	config *config.Config
	s      *ShortlyMux
	l      *shortlyLog.ShortlyLog
}

func NewServer(c *config.Config) *Server {
	server := &Server{}

	// Create shortly log for middleware
	middlewareLog := shortlyLog.ShortlyLog{
		Tag: "HTTP",
	}

	// Create HTTP Server with ShortlyMux
	mux := &ShortlyMux{
		l: &middlewareLog,
	}

	// Initialize ShortlyLog
	slog := new(shortlyLog.ShortlyLog)

	server.config = c
	server.s = mux
	server.l = slog

	return server
}

func (s *Server) Start() error {
	log.Printf("Listening on :%d...", s.config.Server.Port)

	shortlyBase := shortlyDB.ShortlyBase{
		FileName:   "shortly",
		Log:        s.l,
		MemoryPath: s.config.MemoryPath,
	}

	db, err := shortlyBase.InitialDB()

	if err != nil {
		s.l.ZapFatal("Couldn't not connect Shortly DB", err)
	}

	go s.saveToDisk(shortlyBase, db, s.config.DurationOfWriteToDisk)

	shortlyService := shortly.NewShortlyService(db, s.config.LengthOfCode)
	shortlyHandler := shortly.NewShortlyHandler(shortlyService, s.l)

	// For Swagger API Documentation
	// We want to serve our api docs
	fs := http.FileServer(http.Dir("./docs"))
	s.s.Handler(regexp.MustCompile("/docs/"), http.StripPrefix("/docs/", fs))

	// Serve Swagger UI
	s.s.Handler(
		regexp.MustCompile("/api/swagger/.*"),
		httpSwagger.Handler(
			httpSwagger.URL("/docs/swagger.json")))

	s.s.Handler(
		regexp.MustCompile("/api/v1/health"),
		HTTPLogMiddleware(s.s.l, http.HandlerFunc(s.healthCheck)),
	)

	s.s.Handler(
		regexp.MustCompile("/api/v1/create"),
		HTTPLogMiddleware(s.s.l, http.HandlerFunc(shortlyHandler.CreateShortURL)),
	)

	s.s.Handler(
		regexp.MustCompile("api/v1/list"),
		HTTPLogMiddleware(s.s.l, http.HandlerFunc(shortlyHandler.GetShortList)),
	)

	// Short URL End-point.
	s.s.Handler(
		regexp.MustCompile("^/[^/]*$"),
		HTTPLogMiddleware(s.s.l, http.HandlerFunc(shortlyHandler.RedirectURL)),
	)

	err = http.ListenAndServe(
		fmt.Sprintf(":%d", s.config.Server.Port),
		s.s,
	)

	return err
}

// Health Handler
// @Summary Server Health
// @Description It helps server tracking
// @Tags Health
// @Produce text/plain
// @Success 200
// @Router /api/v1/health [get]
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))

	if err != nil {
		s.l.ZapError(shortlyError.ParserError, err)
	}
}

func (s *Server) saveToDisk(shortlyBase shortlyDB.ShortlyBase, db *shortlyDB.DB, durationOfWriteToDisk time.Duration) {
	for {
		time.Sleep(durationOfWriteToDisk)

		localDB, err := shortlyBase.ReadFromFile()

		if err != nil {
			s.l.ZapFatal("couldn't read db", err)
		}

		// Compare memory and file
		if len(localDB.ShortURL) != len(db.ShortURL) {
			// Append new data from file to in-memory
			for _, s := range localDB.ShortURL {
				if db.FindByCode(s.Code) == nil {
					db.ShortURL = append(db.ShortURL, s)
				}
			}

			err := shortlyBase.SaveToFile(db)

			if err != nil {
				shortlyBase.Log.ZapFatal("couldn't save data in disk", err)
				return
			}
		}
	}
}
