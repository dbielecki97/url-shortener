package v1

import (
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"regexp"
	"strings"
)

type ShortenRequest struct {
	URL string `json:"url,omitempty"`
}

const UrlRegex = "^(?:https?:\\/\\/)?(?:[^@\\/\\n]+@)?(?:www\\.)?([^:\\/\\n]+)"

func (r ShortenRequest) Validate() *errs.AppError {
	reg := regexp.MustCompile(UrlRegex)
	matches := reg.Match([]byte(r.URL))
	if !matches {
		return errs.NewValidationError("not a valid url")
	}
	return nil
}

func (r *ShortenRequest) Sanitize() {
	if !strings.HasPrefix(r.URL, "http") {
		r.URL = "https://" + r.URL
	}
}

type ShortenResponse struct {
	Code      string `json:"code,omitempty"`
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}
