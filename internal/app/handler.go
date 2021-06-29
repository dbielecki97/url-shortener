package app

import (
	"encoding/json"
	"fmt"
	"github.com/dbielecki97/url-shortener/internal/api"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

func (s *Server) handleUrlShorten(w http.ResponseWriter, r *http.Request) {
	var req api.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.writeResponse(w, r, http.StatusUnprocessableEntity, fmt.Sprintf("could not unmarshal request body: %v", err))
		return
	}

	response, err := s.service.Shorten(req)
	if err != nil {
		s.log.Errorf("%+v\n", err)
		if isValidationError(err) {
			s.writeResponse(w, r, http.StatusUnprocessableEntity, err.Error())
			return
		}

		s.writeResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeResponse(w, r, http.StatusOK, response)
}

func isValidationError(err error) bool {
	ve, ok := errors.Cause(err).(api.Validator)
	return ok && ve.HasError()
}

func (s *Server) handleUrlExtend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	res, err := s.service.Expand(code)
	if err != nil {
		s.log.Errorf("%+v\n", err)
		if errors.As(err, &domain.NotFoundError{}) {
			s.writeResponse(w, r, http.StatusNotFound, err.Error())
			return
		}

		s.writeResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, res.URL, http.StatusSeeOther)

	s.log.SetReportCaller(false)
	s.log.Infof("Redirected to %v", res.URL)
	s.log.SetReportCaller(true)
}

func (s *Server) handleUrlInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	res, err := s.service.Expand(code)
	if err != nil {
		if errors.As(err, &domain.NotFoundError{}) {
			s.writeResponse(w, r, http.StatusNotFound, err.Error())
			return
		}
		s.log.Errorf("%+x\n", err)
		s.writeResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	s.writeResponse(w, r, http.StatusOK, res)
}

func (s *Server) writeResponse(w http.ResponseWriter, r *http.Request, code int, message interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&message); err != nil {
		s.log.Errorf("%+v\n", errors.Wrap(err, "could not write response body"))
	}
}
