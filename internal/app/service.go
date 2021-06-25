package app

import (
	"github.com/dbielecki97/url-shortener/internal/api"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"time"
)

type Service interface {
	Shorten(api.ShortenRequest) (*api.ShortenResponse, *errs.AppError)
	Expand(code string) (*api.ExpandResponse, *errs.AppError)
}

type DefaultService struct {
	cache domain.ShortUrlRepo
	store domain.ShortUrlRepo
}

func NewDefaultService(cache domain.ShortUrlRepo, store domain.ShortUrlRepo) *DefaultService {
	return &DefaultService{cache: cache, store: store}
}

func (d DefaultService) Shorten(r api.ShortenRequest) (*api.ShortenResponse, *errs.AppError) {
	entry := &domain.ShortURL{
		URL:       r.URL,
		Code:      domain.RandomCode(),
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	entry, err := d.cache.Save(entry)
	if err != nil {
		return nil, err
	}

	entry, err = d.store.Save(entry)
	if err != nil {
		return nil, err
	}

	res := api.ShortenResponse{
		Code:      entry.Code,
		URL:       entry.URL,
		CreatedAt: entry.CreatedAt,
	}

	return &res, nil
}

func (d DefaultService) Expand(code string) (*api.ExpandResponse, *errs.AppError) {
	e, err := d.cache.Find(code)
	if err == nil {
		res := api.ExpandResponse{
			Code:      e.Code,
			URL:       e.URL,
			CreatedAt: e.CreatedAt,
		}

		return &res, nil
	}
	if err.Code != errs.CacheMiss {
		return nil, err
	}

	e, err = d.store.Find(code)
	if err != nil {
		return nil, err
	}

	e, _ = d.cache.Save(e)

	res := api.ExpandResponse{
		Code:      e.Code,
		URL:       e.URL,
		CreatedAt: e.CreatedAt,
	}

	return &res, nil
}
