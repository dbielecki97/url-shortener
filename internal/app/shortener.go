package app

import "github.com/dbielecki97/url-shortener/internal/domain"

var (
	shortener Shortener = &defaultShortener{}
)

type Shortener interface {
	ShortenUrl(url string) *domain.ShortURL
}

type defaultShortener struct{}

func (d defaultShortener) ShortenUrl(url string) *domain.ShortURL {
	return domain.NewShortURL(url)
}
