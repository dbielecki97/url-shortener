package app

import (
	"encoding/json"
	"github.com/dbielecki97/url-shortener/internal/api"
	"github.com/dbielecki97/url-shortener/mocks/app"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var server *Server
var mockService *app.MockService

func setupHandlerTest(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = app.NewMockService(ctrl)
	log := logrus.New()
	server = NewServer(log, mockService)

	return func() {
		defer ctrl.Finish()
	}
}

func Test_handleHealthCheck_should__return_200_with_message(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	request, _ := http.NewRequest(http.MethodGet, "/health", nil)

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}

	type healthCheckResponse struct {
		Message string `json:"message,omitempty"`
	}

	var res healthCheckResponse
	_ = json.NewDecoder(recorder.Body).Decode(&res)

	if res.Message != "OK" {
		t.Errorf("Failed while testing the response body")
	}
}

func Test_handleUrlShorten_should_shorten_url_and_return_200(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	req := api.ShortenRequest{URL: "https://www.google.com"}
	j, _ := json.Marshal(req)
	reader := strings.NewReader(string(j))

	res := api.ShortenInfo{
		Code:      "1231231231",
		URL:       "https://www.google.com",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	mockService.EXPECT().Shorten(req).Return(&res, nil)

	request, _ := http.NewRequest(http.MethodPost, "/", reader)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}
}

func Test_handleUrlShorten_should_shorten_url_and_return_500(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	req := api.ShortenRequest{URL: "https://www.google.com"}
	j, _ := json.Marshal(req)
	reader := strings.NewReader(string(j))

	mockService.EXPECT().Shorten(req).Return(nil, errs.NewUnexpectedError("unexpected database error"))

	request, _ := http.NewRequest(http.MethodPost, "/", reader)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}
}

func Test_handleUrlInfo_should_return_200(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	res := api.ShortenInfo{
		Code:      "123123123a",
		URL:       "https://www.google.com",
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	mockService.EXPECT().Expand("123123123a").Return(&res, nil)

	request, _ := http.NewRequest(http.MethodGet, "/info/123123123a", nil)

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}
}

func Test_handleUrlInfo_should_return_400(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	mockService.EXPECT().Expand("123123123a").Return(nil, errs.NewNotFoundError("invalid code"))

	request, _ := http.NewRequest(http.MethodGet, "/info/123123123a", nil)

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Error("Failed while testing the status code")
	}

	var err errs.AppError
	_ = json.NewDecoder(recorder.Body).Decode(&err)

	if err.Message != "invalid code" {
		t.Error("Failed while testing the response body")
	}
}

func Test_handleUrlInfo_should_return_500(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	mockService.EXPECT().Expand("123123123a").Return(nil, errs.NewUnexpectedError("unexpected database error"))

	request, _ := http.NewRequest(http.MethodGet, "/info/123123123a", nil)

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}

	var err errs.AppError
	_ = json.NewDecoder(recorder.Body).Decode(&err)

	if err.Message != "unexpected database error" {
		t.Error("Failed while testing the response body")
	}
}

func Test_handleUrlExtend_should_return_400(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	mockService.EXPECT().Expand("123123123a").Return(nil, errs.NewNotFoundError("invalid code"))

	request, _ := http.NewRequest(http.MethodGet, "/123123123a", nil)

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Error("Failed while testing the status code")
	}

	var err errs.AppError
	_ = json.NewDecoder(recorder.Body).Decode(&err)

	if err.Message != "invalid code" {
		t.Error("Failed while testing the response body")
	}
}

func Test_handleUrlExtend_should_return_500(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	mockService.EXPECT().Expand("123123123a").Return(nil, errs.NewUnexpectedError("unexpected database error"))

	request, _ := http.NewRequest(http.MethodGet, "/123123123a", nil)

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}

	var err errs.AppError
	_ = json.NewDecoder(recorder.Body).Decode(&err)

	if err.Message != "unexpected database error" {
		t.Error("Failed while testing the response body")
	}
}

func Test_handleUrlExtend_should_return_302(t *testing.T) {
	teardown := setupHandlerTest(t)
	defer teardown()

	url := "https://www.google.com"
	info := api.ShortenInfo{
		Code:      "123123123a",
		URL:       url,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	mockService.EXPECT().Expand("123123123a").Return(&info, nil)

	request, _ := http.NewRequest(http.MethodGet, "/123123123a", nil)

	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusSeeOther {
		t.Error("Failed while testing the status code")
	}

	lh := recorder.Header().Get("Location")

	if lh != url {
		t.Errorf("Failed while testing the response Location header: expected %v, got %v", lh, url)
	}
}
