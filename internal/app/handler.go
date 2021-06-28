package app

import (
	"encoding/json"
	"github.com/dbielecki97/url-shortener/internal/api"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) handleUrlShorten(w http.ResponseWriter, r *http.Request) {
	var req api.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.writeResponse(w, r, http.StatusUnprocessableEntity, "could not unmarshal request body")
		return
	}

	response, appErr := s.service.Shorten(req)
	if appErr != nil {
		s.writeResponse(w, r, appErr.Code, appErr.AsMessage())
		return
	}

	s.writeResponse(w, r, http.StatusOK, response)
}

func (s *Server) handleUrlExtend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	res, err := s.service.Expand(code)
	if err != nil {
		s.writeResponse(w, r, err.Code, err.AsMessage())
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
		s.writeResponse(w, r, err.Code, err.AsMessage())
		return
	}
	s.writeResponse(w, r, http.StatusOK, res)
}

func (s *Server) writeResponse(w http.ResponseWriter, r *http.Request, code int, message interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(&message); err != nil {
		s.log.Errorf("Could not write response body: %v", err)
	}
}
