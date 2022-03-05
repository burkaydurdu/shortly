package shortly

type SaveRequestDTO struct {
	OriginalURL string `json:"original_url"`
}

type SaveResponseDTO struct {
	ShortURL string `json:"short_url"`
}

type ErrResponseDTO struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}
