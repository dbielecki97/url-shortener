package domain

import "time"

type ShortURL struct {
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
	Code      string `json:"code,omitempty"`
}

func NewShortURL(URL string) *ShortURL {
	return &ShortURL{URL: URL, Code: randomCode(), CreatedAt: time.Now().Format(time.RFC3339)}
}
