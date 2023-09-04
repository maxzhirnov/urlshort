package models

import (
	"fmt"
)

type ShortURL struct {
	OriginalURL string `json:"original_url"`
	ID          string `json:"id"`
}

func (u ShortURL) String() string {
	return fmt.Sprintf("%s: %s", u.ID, u.OriginalURL)
}
