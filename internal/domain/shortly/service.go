package shortly

type Service interface {
	RedirectURL(code string) string
}

type shortlyService struct {
	filePath string
}

func (s *shortlyService) RedirectURL(code string) string {
	return "https://github.com/burkaydurdu"
}

func NewShortlyService(filePath string) Service {
	return &shortlyService{
		filePath: filePath,
	}
}
