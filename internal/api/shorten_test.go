package api

import (
	"github.com/pkg/errors"
	"testing"
)

func TestShortenRequest_Validate(t *testing.T) {
	type fields struct {
		URL string
	}
	tests := []struct {
		name   string
		fields fields
		want   error
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
			want:   validationError{err: errors.New("url can't be empty")},
		},
		{
			name:   "Proper URL path with http",
			fields: fields{URL: "http://www.google.com"},
			want:   nil,
		},
		{
			name:   "missing .com",
			fields: fields{URL: "www.google"},
			want:   errors.New("not a valid url"),
		},
		{
			name:   "random sentence as url",
			fields: fields{URL: "not a url at all"},
			want:   errors.New("not a valid url"),
		},
		{
			name:   "valid url without scheme",
			fields: fields{URL: "www.google.com"},
			want:   errors.New("not a valid url"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := ShortenRequest{
				URL: tt.fields.URL,
			}
			got := r.Validate()
			if got == nil {
				if tt.want != got {
					t.Errorf("Validate() = %+v, want %+v", got, tt.want)
				}
			} else if got.Error() != tt.want.Error() {
				t.Errorf("Validate() = %+v, want %+v", got, tt.want)
			}

		})
	}
}
