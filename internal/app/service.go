package app

import (
	"github.com/dbielecki97/url-shortener/internal/api"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/dbielecki97/url-shortener/pkg/logger"
)

//go:generate mockgen -destination=../../mocks/app/mockService.go -package=app github.com/dbielecki97/url-shortener/internal/app Service
type Service interface {
	Shorten(api.ShortenRequest) (*api.ShortenInfo, errs.RestErr)
	Expand(code string) (*api.ShortenInfo, errs.RestErr)
}

type DefaultService struct {
	cache domain.ShortUrlRepo
	store domain.ShortUrlRepo
}

func NewDefaultService(c domain.ShortUrlRepo, s domain.ShortUrlRepo) *DefaultService {
	return &DefaultService{cache: c, store: s}
}

func (d DefaultService) Shorten(r api.ShortenRequest) (*api.ShortenInfo, errs.RestErr) {
	err := r.Validate()
	if err != nil {
		return nil, err
	}

	entry := shortener.ShortenUrl(r.URL)

	entry, restErr := d.cache.Save(entry)
	if restErr != nil {
		return nil, restErr
	}

	entry, restErr = d.store.Save(entry)
	if restErr != nil {
		return nil, restErr
	}

	res := api.ShortenInfo{
		Code:      entry.Code,
		URL:       entry.URL,
		CreatedAt: entry.CreatedAt,
	}

	return &res, nil
}

func (d DefaultService) Expand(code string) (*api.ShortenInfo, errs.RestErr) {
	su, err := d.cache.Find(code)
	if err != nil {
		return nil, err
	}

	if su.URL != "" {
		res := api.ShortenInfo{
			Code:      su.Code,
			URL:       su.URL,
			CreatedAt: su.CreatedAt,
		}

		return &res, nil
	}

	su, err = d.store.Find(code)
	if err != nil {
		return nil, err
	}

	_, restErr := d.cache.Save(su)
	if restErr != nil {
		logger.Error("could not save to cache after reading from client: %v", restErr)
	}

	res := api.ShortenInfo{
		Code:      su.Code,
		URL:       su.URL,
		CreatedAt: su.CreatedAt,
	}

	return &res, nil
}
