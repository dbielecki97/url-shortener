package url

import (
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Entry struct {
	URL       string `json:"url,omitempty"`
	CreatedAt string `json:"created_at,omitempty" db:"created_at"`
	Code      string `json:"code,omitempty"`
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandomCode() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type Repository interface {
	Save(e *Entry) (*Entry, *errs.AppError)
	Find(code string) (*Entry, *errs.AppError)
}
