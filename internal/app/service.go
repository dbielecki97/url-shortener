package app

import (
	"github.com/dbielecki97/url-shortener/internal/api"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -destination=../../mocks/app/mockService.go -package=app github.com/dbielecki97/url-shortener/internal/app Service
type Service interface {
	Shorten(api.ShortenRequest) (*api.ShortenInfo, error)
	Expand(code string) (*api.ShortenInfo, error)
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

func (d DefaultService) Shorten(r api.ShortenRequest) (*api.ShortenInfo, error) {
	err := r.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "validation failed")
	}

	entry := d.shortener.ShortenUrl(r.URL)

	entry, err = d.cache.Save(entry)
	if err != nil {
		return nil, errors.Wrap(err, "could not save to cache")
	}

	entry, err = d.store.Save(entry)
	if err != nil {
		return nil, errors.Wrap(err, "could not save to store")
	}

	res := api.ShortenInfo{
		Code:      entry.Code,
		URL:       entry.URL,
		CreatedAt: entry.CreatedAt,
	}

	return &res, nil
}

func IsCacheMiss(err error) bool {
	return errors.As(err, &domain.NotFoundError{})
}

func (d DefaultService) Expand(code string) (*api.ShortenInfo, error) {
	info, err := d.cache.Find(code)
	if err == nil {
		res := api.ShortenInfo{
			Code:      info.Code,
			URL:       info.URL,
			CreatedAt: info.CreatedAt,
		}

		return &res, nil
	}

	if !IsCacheMiss(err) {
		return nil, errors.Wrap(err, "could not find entity")
	}

	info, err = d.store.Find(code)
	if err != nil {
		return nil, errors.Wrap(err, "could not find entity")
	}

	_, err = d.cache.Save(info)
	if err != nil {
		d.log.Errorf("%+v\n", errors.Wrap(err, "could not save entity"))
	}

	res := api.ShortenInfo{
		Code:      info.Code,
		URL:       info.URL,
		CreatedAt: info.CreatedAt,
	}

	return &res, nil
}
