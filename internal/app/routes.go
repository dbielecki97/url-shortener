package app

import (
	"github.com/dbielecki97/url-shortener/pkg/resp"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) routes() {
	s.router = mux.NewRouter()

	s.router.HandleFunc("/health", s.handleHealthCheck)
	s.router.HandleFunc("/", s.handleUrlShorten).Methods(http.MethodPost)
	s.router.HandleFunc("/{code:[a-zA-Z0-9]{10}}", s.handleUrlExtend).Methods(http.MethodGet)
	s.router.HandleFunc("/info/{code:[a-zA-Z0-9]{10}}", s.handleUrlInfo).Methods(http.MethodGet)

}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	resp.JSON(w, http.StatusOK, struct {
		Message string `json:"message,omitempty"`
	}{Message: "OK"})
}
