package domain

import (
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"math/rand"
	"time"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

type ShortURL struct {
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
	Code      string `json:"code,omitempty"`
}

func NewShortURL(URL string) *ShortURL {
	return &ShortURL{URL: URL, Code: randomCode(), CreatedAt: time.Now().Format(time.RFC3339)}
}

func randomCode() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//go:generate mockgen -destination=../../mocks/domain/mockRepository.go -package=domain github.com/dbielecki97/url-shortener/internal/domain ShortUrlRepo
type ShortUrlRepo interface {
	Save(e *ShortURL) (*ShortURL, errs.RestErr)
	Find(code string) (*ShortURL, errs.RestErr)
}
