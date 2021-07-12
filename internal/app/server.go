package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	service Service
	router  *mux.Router
}

func NewServer(service Service) *Server {
	s := &Server{service: service}
	s.routes()
	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
