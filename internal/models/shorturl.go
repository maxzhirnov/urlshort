package models

import (
	"fmt"
)

type ShortURL struct {
	OriginalURL string `json:"original_url"`
	ID          string `json:"id"`
	UUID        string `json:"uuid"`
	DeletedFlag bool   `json:"deleted_flag"`
}

func (u ShortURL) String() string {
	return fmt.Sprintf("%s: %s", u.ID, u.OriginalURL)
}
