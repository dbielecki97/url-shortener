package domain

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -destination=../../mocks/domain/mockRepository.go -package=domain github.com/dbielecki97/url-shortener/internal/domain ShortUrlRepo
type ShortUrlRepo interface {
	Save(e *ShortURL) (*ShortURL, error)
	Find(code string) (*ShortURL, error)
}

type NotFoundError struct {
	Err error
}

func (s NotFoundError) Error() string {
	return s.Err.Error()
}
