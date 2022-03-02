package server

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/burkaydurdu/shortly/internal/domain/shortly"

	"github.com/burkaydurdu/shortly/config"
	shortlyError "github.com/burkaydurdu/shortly/pkg/error"
)

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type ShortlyMux struct {
	routes []*route
	l      *ShortlyLog
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
	l      *ShortlyLog
}

func NewServer(c *config.Config) *Server {
	server := &Server{}

	// Create shortly log for middleware
	middlewareLog := ShortlyLog{
		Tag: "HTTP",
	}

	// Create HTTP Server with ShortlyMux
	mux := &ShortlyMux{
		l: &middlewareLog,
	}

	// Initialize ShortlyLog
	shortlyLog := new(ShortlyLog)

	server.config = c
	server.s = mux
	server.l = shortlyLog

	return server
}

func (s *Server) Start() error {
	log.Printf("Listening on :%d...", s.config.Server.Port)

	shortlyService := shortly.NewShortlyService("./data")
	shortlyHandler := shortly.NewShortlyHandler(shortlyService)

	s.s.Handler(
		regexp.MustCompile("/api/health"),
		HTTPLogMiddleware(s.s.l, http.HandlerFunc(s.healthCheck)),
	)

	s.s.Handler(
		regexp.MustCompile("^/[^/]*$"),
		HTTPLogMiddleware(s.s.l, http.HandlerFunc(shortlyHandler.RedirectURL)),
	)

	err := http.ListenAndServe(
		fmt.Sprintf(":%d", s.config.Server.Port),
		s.s,
	)

	return err
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("OK"))

	if err != nil {
		s.l.ZapError(shortlyError.ParserError, err)
	}
}
