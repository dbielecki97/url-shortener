package app

import (
	"github.com/dbielecki97/url-shortener/internal/api"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=../../mocks/app/mockService.go -package=app github.com/dbielecki97/url-shortener/internal/app Service
type Service interface {
	Shorten(api.ShortenRequest) (*api.ShortenInfo, *errs.AppError)
	Expand(code string) (*api.ShortenInfo, *errs.AppError)
}

type DefaultService struct {
	cache     domain.ShortUrlRepo
	store     domain.ShortUrlRepo
	shortener Shortener
	log       *logrus.Logger
}

func NewDefaultService(c domain.ShortUrlRepo, s domain.ShortUrlRepo, log *logrus.Logger, st Shortener) *DefaultService {
	return &DefaultService{cache: c, store: s, log: log, shortener: st}
}

func (d DefaultService) Shorten(r api.ShortenRequest) (*api.ShortenInfo, *errs.AppError) {
	err := r.Validate()
	if err != nil {
		return nil, err
	}

	entry := d.shortener.ShortenUrl(r.URL)

	entry, err = d.cache.Save(entry)
	if err != nil {
		return nil, err
	}

	entry, err = d.store.Save(entry)
	if err != nil {
		return nil, err
	}

	res := api.ShortenInfo{
		Code:      entry.Code,
		URL:       entry.URL,
		CreatedAt: entry.CreatedAt,
	}

	return &res, nil
}

func (d DefaultService) Expand(code string) (*api.ShortenInfo, *errs.AppError) {
	e, err := d.cache.Find(code)
	if err == nil {
		res := api.ShortenInfo{
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

	_, err = d.cache.Save(e)
	if err != nil {
		d.log.Errorf("Could not save to cache after reading from store: %v", err)
	}

	res := api.ShortenInfo{
		Code:      e.Code,
		URL:       e.URL,
		CreatedAt: e.CreatedAt,
	}

	return &res, nil
}
