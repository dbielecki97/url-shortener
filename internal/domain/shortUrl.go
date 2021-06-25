package domain

type ShortURL struct {
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
	Code      string `json:"code,omitempty"`
}
