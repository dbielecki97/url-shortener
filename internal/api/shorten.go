package api

import (
	"github.com/pkg/errors"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url,omitempty"`
}

type Validator interface {
	HasError() bool
}

type validationError struct {
	err error
}

func (v validationError) HasError() bool {
	return true
}

func (v validationError) Error() string {
	return v.err.Error()
}

func (r ShortenRequest) Validate() error {
	if r.URL == "" {
		return validationError{err: errors.New("url can't be empty")}
	}

	_, err := url.ParseRequestURI(r.URL)
	if err != nil {
		return validationError{err: errors.New("not a valid url")}
	}

	return nil
}

type ShortenInfo struct {
	Code      string `json:"code,omitempty"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
