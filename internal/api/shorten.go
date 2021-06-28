package api

import (
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url,omitempty"`
}

func (r ShortenRequest) Validate() *errs.AppError {
	if r.URL == "" {
		return errs.NewValidationError("url can't be empty")
	}

	_, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return errs.NewValidationError("not a valid url")
	}

	return nil
}

type ShortenInfo struct {
	Code      string `json:"code,omitempty"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
