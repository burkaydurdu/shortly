package shortly

import (
	"encoding/json"
	"net/http"

	"github.com/burkaydurdu/shortly/pkg/log"

	shortlyError "github.com/burkaydurdu/shortly/pkg/error"
)

type HTTPHandler interface {
	RedirectURL(w http.ResponseWriter, r *http.Request)
	SaveShortURL(w http.ResponseWriter, r *http.Request)
	GetShortList(w http.ResponseWriter, r *http.Request)
}

type shortlyHandler struct {
	s Service
	l *log.ShortlyLog
}

func (s *shortlyHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	// Remove first character `` / ``
	code := r.URL.Path[1:]

	redirectURL, err := s.s.RedirectURL(code)

	if err != nil {
		http.NotFound(w, r)

		return
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (s *shortlyHandler) SaveShortURL(w http.ResponseWriter, r *http.Request) {
	var requestBody SaveRequestDTO

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		httpResponse(w, []byte(shortlyError.ParserError), s.l)

		return
	}

	responseBody := s.s.SaveShortURL(r.Host, &requestBody)

	responseByteBody, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	httpResponse(w, responseByteBody, s.l)
}

func (s *shortlyHandler) GetShortList(w http.ResponseWriter, _ *http.Request) {
	shortlyURL := s.s.GetShortList()

	responseByteBody, _ := json.Marshal(shortlyURL)

	w.Header().Set("Content-Type", "application/json")
	httpResponse(w, responseByteBody, s.l)
}

func httpResponse(w http.ResponseWriter, response []byte, l *log.ShortlyLog) {
	_, err := w.Write(response)
	if err != nil {
		l.ZapFatal(shortlyError.ResponseError, err)
	}
}

func NewShortlyHandler(s Service, l *log.ShortlyLog) HTTPHandler {
	return &shortlyHandler{
		s: s,
		l: l,
	}
}
