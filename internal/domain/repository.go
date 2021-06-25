package domain

import "github.com/dbielecki97/url-shortener/pkg/errs"

type ShortUrlRepo interface {
	Save(e *ShortURL) (*ShortURL, *errs.AppError)
	Find(code string) (*ShortURL, *errs.AppError)
}
