package shortly

type SaveRequestDTO struct {
	OriginalUrl string `json:"original_url"`
}

type SaveResponseDTO struct {
	ShortURL string `json:"short_url"`
}
