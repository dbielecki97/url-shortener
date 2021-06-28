package domain

import "github.com/dbielecki97/url-shortener/pkg/errs"
import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -destination=../../mocks/domain/mockRepository.go -package=domain github.com/dbielecki97/url-shortener/internal/domain ShortUrlRepo
type ShortUrlRepo interface {
	Save(e *ShortURL) (*ShortURL, *errs.AppError)
	Find(code string) (*ShortURL, *errs.AppError)
}
