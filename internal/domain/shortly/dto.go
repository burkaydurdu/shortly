package shortly

type SaveRequestDTO struct {
	OriginalURL string `json:"original_url"`
}

type SaveResponseDTO struct {
	ShortURL string `json:"short_url"`
}
