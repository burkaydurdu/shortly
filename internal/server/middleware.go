package server

import (
	"fmt"
	"net/http"
)

func HTTPLogMiddleware(l *ShortlyLog, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Zap(fmt.Sprintf("%s %s", r.Method, r.URL.Path))

		next.ServeHTTP(w, r)
	})
}
