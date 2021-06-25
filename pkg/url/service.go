package url

import (
	"github.com/dbielecki97/url-shortener/api/v1"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"time"
)

type Service interface {
	Shorten(v1.ShortenRequest) (*v1.ShortenResponse, *errs.AppError)
	Expand(code string) (*v1.ExpandResponse, *errs.AppError)
}

type DefaultService struct {
	cache Repository
	store Repository
}

func NewService(cache Repository, store Repository) *DefaultService {
	return &DefaultService{cache: cache, store: store}
}

func (d DefaultService) Shorten(r v1.ShortenRequest) (*v1.ShortenResponse, *errs.AppError) {
	entry := &Entry{
		URL:       r.URL,
		Code:      RandomCode(),
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

	res := v1.ShortenResponse{
		Code:      entry.Code,
		URL:       entry.URL,
		CreatedAt: entry.CreatedAt,
	}

	return &res, nil
}

func (d DefaultService) Expand(code string) (*v1.ExpandResponse, *errs.AppError) {
	e, err := d.cache.Find(code)
	if err == nil {
		res := v1.ExpandResponse{
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

	res := v1.ExpandResponse{
		Code:      e.Code,
		URL:       e.URL,
		CreatedAt: e.CreatedAt,
	}

	return &res, nil
}
