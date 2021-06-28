package api

import (
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"reflect"
	"testing"
)

func TestShortenRequest_Validate(t *testing.T) {
	type fields struct {
		URL string
	}
	tests := []struct {
		name   string
		fields fields
		want   *errs.AppError
	}{
		{
			name:   "Proper URL path with https",
			fields: fields{URL: "https://www.google.com"},
			want:   nil,
		},
		{
			name:   "Proper URL path with only https",
			fields: fields{URL: "https://"},
			want:   nil,
		},
		{
			name:   "Proper URL path with only https",
			fields: fields{URL: ""},
			want:   errs.NewValidationError("url can't be empty"),
		},
		{
			name:   "Proper URL path with http",
			fields: fields{URL: "http://www.google.com"},
			want:   nil,
		},
		{
			name:   "missing .com",
			fields: fields{URL: "www.google"},
			want:   errs.NewValidationError("not a valid url"),
		},
		{
			name:   "random sentence as url",
			fields: fields{URL: "not a url at all"},
			want:   errs.NewValidationError("not a valid url"),
		},
		{
			name:   "valid url without scheme",
			fields: fields{URL: "www.google.com"},
			want:   errs.NewValidationError("not a valid url"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ShortenRequest{
				URL: tt.fields.URL,
			}
			if got := r.Validate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
