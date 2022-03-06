package shortly

import (
	"encoding/json"
	"net/http"

	"github.com/burkaydurdu/shortly/internal/db"

	"github.com/burkaydurdu/shortly/pkg/util"

	"github.com/burkaydurdu/shortly/pkg/log"

	shortlyError "github.com/burkaydurdu/shortly/pkg/error"
)

type HTTPHandler interface {
	RedirectURL(w http.ResponseWriter, r *http.Request)
	CreateShortURL(w http.ResponseWriter, r *http.Request)
	GetShortList(w http.ResponseWriter, r *http.Request)
}

type shortlyHandler struct {
	s Service
	l *log.ShortlyLog
}

// RedirectURL Handler
// @Summary Short URL
// @Description It redirects from short URL to original URL
// @Tags Redirect
// @Success 200 "Redirect URL"
// @Failure 404 {string} string "Not found"
// @Param code path string true "Shortly Code"
// @Router /{code} [get]
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

// CreateShortURL
// @Summary Generate short URL
// @Description It gives short URL
// @Tags Create Short URL
// @Accept json
// @Produce json
// @Param original_url body shortly.SaveRequestDTO true "Original URL is required"
// @Success 200 {object} shortly.SaveResponseDTO
// @failure 400 {object} shortly.ErrResponseDTO object "Request body is not valid"
// @failure 409 {object} shortly.ErrResponseDTO object "URL is not valid"
// @Router /api/v1/create [post]
func (s *shortlyHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var requestBody SaveRequestDTO

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&requestBody)

	if err != nil {
		body := responseErrorBody(err, shortlyError.ParserError, shortlyError.ParserErrCode)

		shortlyResponse(w, body, s.l, http.StatusBadRequest)

		return
	}

	if err = util.IsURL(requestBody.OriginalURL); err != nil {
		body := responseErrorBody(err, shortlyError.InvalidURLError, shortlyError.InvalidParams)

		shortlyResponse(w, body, s.l, http.StatusConflict)

		return
	}

	shortlyURL := s.s.CreateShortURL(r.Host, &requestBody)
	responseBody, _ := json.Marshal(shortlyURL)

	shortlyResponse(w, responseBody, s.l, http.StatusOK)
}

// GetShortList
// @Summary Get All Shortly List
// @Description It gives all shortly data
// @Tags All Shortly List
// @Success 200 {array} db.Shortly
// @Router /api/v1/list [get]
func (s *shortlyHandler) GetShortList(w http.ResponseWriter, _ *http.Request) {
	shortlyURL := s.s.GetShortList()

	// when the list is empty it returns null therefore we should create this object.
	if len(shortlyURL) == 0 {
		shortlyURL = make([]db.Shortly, 0)
	}

	responseBody, _ := json.Marshal(shortlyURL)

	shortlyResponse(w, responseBody, s.l, http.StatusOK)
}

func shortlyResponse(w http.ResponseWriter, response []byte, l *log.ShortlyLog, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err := w.Write(response)

	if err != nil {
		l.ZapFatal(shortlyError.ResponseError, err)
	}
}

func responseErrorBody(err error, message string, code int) []byte {
	errResponse := ErrResponseDTO{
		Error:   err.Error(),
		Message: message,
		Code:    code,
	}

	result, _ := json.Marshal(errResponse)

	return result
}

func NewShortlyHandler(s Service, l *log.ShortlyLog) HTTPHandler {
	return &shortlyHandler{
		s: s,
		l: l,
	}
}
