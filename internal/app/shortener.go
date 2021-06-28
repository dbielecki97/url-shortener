package app

import "github.com/dbielecki97/url-shortener/internal/domain"

type Shortener interface {
	ShortenUrl(url string) *domain.ShortURL
}

type DefaultShortener struct{}

func (d DefaultShortener) ShortenUrl(url string) *domain.ShortURL {
	return domain.NewShortURL(url)
}
