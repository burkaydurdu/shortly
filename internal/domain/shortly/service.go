package shortly

import (
	"fmt"

	"github.com/burkaydurdu/shortly/internal/db"
	shortlyError "github.com/burkaydurdu/shortly/pkg/error"
	"github.com/burkaydurdu/shortly/pkg/util"
)

type Service interface {
	RedirectURL(code string) (string, error)
	SaveShortURL(host string, requestBody *SaveRequestDTO) *SaveResponseDTO
	GetShortList() []db.Shortly
}

type shortlyService struct {
	db           *db.DB
	lengthOfCode int
}

func (s *shortlyService) RedirectURL(code string) (string, error) {
	for i, sURL := range s.db.ShortURL {
		if code == sURL.Code {
			// Increment visit count
			s.db.ShortURL[i].VisitCount++
			return sURL.OriginalURL, nil
		}
	}

	return "", shortlyError.ErrAddressNotFound
}

func (s *shortlyService) SaveShortURL(host string, requestBody *SaveRequestDTO) *SaveResponseDTO {
	code := generateShortlyCode(s.db.ShortURL, s.lengthOfCode)

	shortly := db.Shortly{
		OriginalURL: requestBody.OriginalURL,
		Code:        code,
		VisitCount:  0,
		ShortURL:    fmt.Sprintf("http://%s/%s", host, code),
	}

	s.db.ShortURL = append(s.db.ShortURL, shortly)

	var responseBody = &SaveResponseDTO{
		ShortURL: shortly.ShortURL,
	}

	return responseBody
}

func (s *shortlyService) GetShortList() []db.Shortly {
	return s.db.ShortURL
}

/*
 * Shortly Code Generator
 * This Method uses RandShortlyCode. It generates six random letters.
 * If it generates same code, it will run again
 * We can generate 54 x 54 x 54 x 54 x 54 x 54 different values
 */
func generateShortlyCode(list []db.Shortly, lengthOfCode int) string {
	code := util.RandShortlyCode(lengthOfCode)

	for _, s := range list {
		if code == s.Code {
			return generateShortlyCode(list, lengthOfCode)
		}
	}

	return code
}

func NewShortlyService(shortlyDB *db.DB, lengthOfCode int) Service {
	return &shortlyService{
		db:           shortlyDB,
		lengthOfCode: lengthOfCode,
	}
}
