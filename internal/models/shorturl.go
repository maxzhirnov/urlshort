package models

type ShortURL struct {
	OriginalURL string `json:"original_url"`
	ID          string `json:"id"`
}
