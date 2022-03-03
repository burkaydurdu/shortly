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
	db *db.DB
}

func (s *shortlyService) RedirectURL(code string) (string, error) {
	for i, sURL := range s.db.ShortUrl {
		if code == sURL.Code {
			// Increment visit count
			s.db.ShortUrl[i].VisitCount += 1
			return sURL.OriginalURL, nil
		}
	}

	return "", shortlyError.AddressNotFoundErr
}

func (s *shortlyService) SaveShortURL(host string, requestBody *SaveRequestDTO) *SaveResponseDTO {
	code := generateShortlyCode(s.db.ShortUrl)

	shortly := db.Shortly{
		OriginalURL: requestBody.OriginalUrl,
		Code:        code,
		VisitCount:  0,
		ShortURL:    fmt.Sprintf("http://%s/%s", host, code),
	}

	s.db.ShortUrl = append(s.db.ShortUrl, shortly)

	var responseBody = &SaveResponseDTO{
		ShortURL: shortly.ShortURL,
	}

	return responseBody
}

func (s *shortlyService) GetShortList() []db.Shortly {
	return s.db.ShortUrl
}

/*
 * Shortly Code Generator
 * This Method uses RandShortlyCode. It generates six random letters.
 * If it generates same code, it will run again
 * We can generate 54 x 54 x 54 x 54 x 54 x 54 different values
 */
func generateShortlyCode(list []db.Shortly) string {
	code := util.RandShortlyCode(6)

	for _, s := range list {
		if code == s.Code {
			return generateShortlyCode(list)
		}
	}

	return code
}

func NewShortlyService(db *db.DB) Service {
	return &shortlyService{
		db: db,
	}
}
