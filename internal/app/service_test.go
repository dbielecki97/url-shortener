package app

import (
	"github.com/dbielecki97/url-shortener/internal/api"
	realDomain "github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/dbielecki97/url-shortener/mocks/domain"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"net/http"
	"testing"
)

var mockCacheRepo *domain.MockShortUrlRepo
var mockStoreRepo *domain.MockShortUrlRepo
var service Service

type mockShortener struct {
}

func (m mockShortener) ShortenUrl(url string) *realDomain.ShortURL {
	return &realDomain.ShortURL{
		URL:       url,
		Code:      "123123123a",
		CreatedAt: "2006-01-02T15:04:05Z07:00",
	}
}

func setupServiceTest(t *testing.T) func() {
	ctrl1 := gomock.NewController(t)
	ctrl2 := gomock.NewController(t)
	mockCacheRepo = domain.NewMockShortUrlRepo(ctrl1)
	mockStoreRepo = domain.NewMockShortUrlRepo(ctrl2)
	service = NewDefaultService(mockCacheRepo, mockStoreRepo, logrus.New(), mockShortener{})

	return func() {
		service = nil
		defer ctrl1.Finish()
		defer ctrl2.Finish()
	}
}

func Test_Shorten_should_receive_error_from_validate_method_invalid_url(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	request := api.ShortenRequest{URL: "invalid url"}

	_, appError := service.Shorten(request)

	if appError == nil {
		t.Error("failed while testing the new account validation")
	}
}

func Test_Shorten_should_receive_error_from_validate_method_not_schema(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	request := api.ShortenRequest{URL: "www.google.com"}

	_, appError := service.Shorten(request)

	if appError.Message != "not a valid url" {
		t.Error("failed while testing the new account validation")
	}
}

func Test_Shorten_should_receive_error_from_validate_method_empty_url(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	request := api.ShortenRequest{URL: ""}

	_, appError := service.Shorten(request)

	if appError.Message != "url can't be empty" {
		t.Error("failed while testing the new account validation")
	}
}

func Test_Shorten_should_receive_error_from_repository_when_cache_can_not_save_shortened_url(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	req := api.ShortenRequest{URL: "https://www.google.com"}

	a := realDomain.ShortURL{
		URL:       "https://www.google.com",
		Code:      "123123123a",
		CreatedAt: "2006-01-02T15:04:05Z07:00",
	}

	mockCacheRepo.EXPECT().Save(&a).Return(nil, errs.NewUnexpectedError("unexpected database error"))
	_, appError := service.Shorten(req)

	if appError.Message != "unexpected database error" {
		t.Error("Test failed while testing for unexpected errors")
	}
}

func Test_Shorten_should_receive_error_from_repository_when_store_can_not_save_shortened_url(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	req := api.ShortenRequest{URL: "https://www.google.com"}

	a := realDomain.ShortURL{
		URL:       "https://www.google.com",
		Code:      "123123123a",
		CreatedAt: "2006-01-02T15:04:05Z07:00",
	}

	mockCacheRepo.EXPECT().Save(&a).Return(&a, nil)
	mockStoreRepo.EXPECT().Save(&a).Return(nil, errs.NewUnexpectedError("unexpected database error"))

	_, appError := service.Shorten(req)

	if appError.Message != "unexpected database error" {
		t.Error("Test failed while testing for unexpected errors")
	}
}

func Test_Shorten_should_receive_error_from_repository_when_shortened_url_was_saved(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	req := api.ShortenRequest{URL: "https://www.google.com"}

	a := realDomain.ShortURL{
		URL:       "https://www.google.com",
		Code:      "123123123a",
		CreatedAt: "2006-01-02T15:04:05Z07:00",
	}

	mockCacheRepo.EXPECT().Save(&a).Return(&a, nil)
	mockStoreRepo.EXPECT().Save(&a).Return(&a, nil)

	res, _ := service.Shorten(req)

	if res.URL != a.URL {
		t.Errorf("Test failed while testing response url, expected %v, got %v", res.URL, a.URL)
	}

	if res.Code != a.Code {
		t.Errorf("Test failed while testing response code, expected %v, got %v", res.Code, a.Code)
	}

	if res.CreatedAt != a.CreatedAt {
		t.Errorf("Test failed while testing response createAt, expected %v, got %v", res.CreatedAt, a.CreatedAt)
	}
}

func Test_Expand_should_return_shortened_url_read_from_cache(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	a := realDomain.ShortURL{
		URL:       "https://www.google.com",
		Code:      "123123123a",
		CreatedAt: "2006-01-02T15:04:05Z07:00",
	}

	mockCacheRepo.EXPECT().Find("123123123a").Return(&a, nil)

	res, _ := service.Expand("123123123a")

	if res.URL != a.URL {
		t.Errorf("Test failed while testing response url, expected %v, got %v", res.URL, a.URL)
	}

	if res.Code != a.Code {
		t.Errorf("Test failed while testing response code, expected %v, got %v", res.Code, a.Code)
	}

	if res.CreatedAt != a.CreatedAt {
		t.Errorf("Test failed while testing response createAt, expected %v, got %v", res.CreatedAt, a.CreatedAt)
	}
}

func Test_Expand_should_return_shortened_url_read_from_store_and_save_to_cache(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	a := realDomain.ShortURL{
		URL:       "https://www.google.com",
		Code:      "123123123a",
		CreatedAt: "2006-01-02T15:04:05Z07:00",
	}

	mockCacheRepo.EXPECT().Find("123123123a").Return(nil, errs.NewCacheMissError())
	mockStoreRepo.EXPECT().Find("123123123a").Return(&a, nil)
	mockCacheRepo.EXPECT().Save(&a).Return(&a, nil)

	res, _ := service.Expand("123123123a")

	if res.URL != a.URL {
		t.Errorf("Test failed while testing response url, expected %v, got %v", res.URL, a.URL)
	}

	if res.Code != a.Code {
		t.Errorf("Test failed while testing response code, expected %v, got %v", res.Code, a.Code)
	}

	if res.CreatedAt != a.CreatedAt {
		t.Errorf("Test failed while testing response createAt, expected %v, got %v", res.CreatedAt, a.CreatedAt)
	}
}

func Test_Expand_should_return_shortened_url_read_from_store_with_error_saving_to_cache(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	a := realDomain.ShortURL{
		URL:       "https://www.google.com",
		Code:      "123123123a",
		CreatedAt: "2006-01-02T15:04:05Z07:00",
	}

	mockCacheRepo.EXPECT().Find("123123123a").Return(nil, errs.NewCacheMissError())
	mockStoreRepo.EXPECT().Find("123123123a").Return(&a, nil)
	mockCacheRepo.EXPECT().Save(&a).Return(nil, errs.NewUnexpectedError("unexpected database error"))

	res, _ := service.Expand("123123123a")

	if res.URL != a.URL {
		t.Errorf("Test failed while testing response url, expected %v, got %v", res.URL, a.URL)
	}

	if res.Code != a.Code {
		t.Errorf("Test failed while testing response code, expected %v, got %v", res.Code, a.Code)
	}

	if res.CreatedAt != a.CreatedAt {
		t.Errorf("Test failed while testing response createAt, expected %v, got %v", res.CreatedAt, a.CreatedAt)
	}
}

func Test_Expand_should_not_found_shortened_url_and_return_not_found_error(t *testing.T) {
	teardown := setupServiceTest(t)
	defer teardown()

	mockCacheRepo.EXPECT().Find("123123123b").Return(nil, errs.NewCacheMissError())
	mockStoreRepo.EXPECT().Find("123123123b").Return(nil, errs.NewNotFoundError("incorrect code"))

	_, err := service.Expand("123123123b")

	if err.Code != http.StatusNotFound {
		t.Errorf("Failed when testing for notFound status code")
	}
}
