package shortly

import (
	"net/http"
)

type HTTPHandler interface {
	RedirectURL(w http.ResponseWriter, r *http.Request)
}

type shortlyHandler struct {
	s Service
}

func (s shortlyHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	// Remove first character `` / ``
	code := r.URL.Path[1:]

	redirectURL := s.s.RedirectURL(code)

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func NewShortlyHandler(s Service) HTTPHandler {
	return &shortlyHandler{
		s: s,
	}
}
