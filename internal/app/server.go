package app

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	log     *logrus.Logger
	service Service
	router  *mux.Router
}

func NewServer(log *logrus.Logger, service Service) *Server {
	s := &Server{log: log, service: service}
	s.routes()
	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
