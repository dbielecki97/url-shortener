package app

import (
	"net/http"
)

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
