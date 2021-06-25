package server

import (
	"encoding/json"
	"github.com/dbielecki97/url-shortener/api/v1"
	"github.com/dbielecki97/url-shortener/pkg/url"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	log     *logrus.Logger
	service url.Service
	router  *mux.Router
}

func New(log *logrus.Logger, service url.Service) *Server {
	s := &Server{log: log, service: service}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router = mux.NewRouter()

	s.router.HandleFunc("/health", s.handleHealthCheck)
	s.router.HandleFunc("/", s.handleUrlShorten).Methods(http.MethodPost)
	s.router.HandleFunc("/{code:[a-zA-Z0-9]{10}}", s.handleUrlExtend).Methods(http.MethodGet)
	s.router.HandleFunc("/info/{code:[a-zA-Z0-9]{10}}", s.handleUrlInfo).Methods(http.MethodGet)

	s.router.Use(s.middlewareLogging())
}

func (s *Server) middlewareLogging() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.log.SetReportCaller(false)
			s.log.Infof(" %v %v", r.Method, r.RequestURI)

			s.log.SetReportCaller(true)
			next.ServeHTTP(w, r)
		})
	}
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	s.writeResponse(w, r, http.StatusOK, v1.HealthCheckResponse{Message: "OK"})
}

func (s *Server) handleUrlShorten(w http.ResponseWriter, r *http.Request) {
	var req v1.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.writeResponse(w, r, http.StatusUnprocessableEntity, "could not unmarshal request body")
		return
	}

	appErr := req.Validate()
	if appErr != nil {
		s.writeResponse(w, r, appErr.Code, appErr.AsMessage())
		return
	}

	req.Sanitize()

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
