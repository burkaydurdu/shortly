package shortly

import (
	"encoding/json"
	"net/http"

	shortlyError "github.com/burkaydurdu/shortly/pkg/error"
)

type HTTPHandler interface {
	RedirectURL(w http.ResponseWriter, r *http.Request)
	SaveShortURL(w http.ResponseWriter, r *http.Request)
	GetShortList(w http.ResponseWriter, r *http.Request)
}

type shortlyHandler struct {
	s Service
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
		w.Write([]byte(shortlyError.ParserError))
		return
	}

	responseBody := s.s.SaveShortURL(r.Host, &requestBody)

	responseByteBody, _ := json.Marshal(responseBody)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseByteBody)
}

func (s *shortlyHandler) GetShortList(w http.ResponseWriter, _ *http.Request) {
	shortlyUrl := s.s.GetShortList()

	responseByteBody, _ := json.Marshal(shortlyUrl)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseByteBody)
}

func NewShortlyHandler(s Service) HTTPHandler {
	return &shortlyHandler{
		s: s,
	}
}
