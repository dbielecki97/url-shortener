package app

import (
	"encoding/json"
	"github.com/dbielecki97/url-shortener/internal/api"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/dbielecki97/url-shortener/pkg/resp"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) handleUrlShorten(w http.ResponseWriter, r *http.Request) {
	var req api.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		resp.Error(w, errs.NewErr("could not unmarshal request body", http.StatusUnprocessableEntity, "unprocessable_entity"))
		return
	}

	response, restErr := s.service.Shorten(req)
	if restErr != nil {
		resp.Error(w, restErr)
		return
	}

	resp.JSON(w, http.StatusOK, response)
}

func (s *Server) handleUrlExtend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	res, err := s.service.Expand(code)
	if err != nil {
		resp.Error(w, err)
		return
	}

	http.Redirect(w, r, res.URL, http.StatusSeeOther)
}

func (s *Server) handleUrlInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["code"]

	res, err := s.service.Expand(code)
	if err != nil {
		resp.Error(w, err)
		return
	}

	resp.JSON(w, http.StatusOK, res)
}
