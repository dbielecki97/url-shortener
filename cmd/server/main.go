package main

import (
	"github.com/dbielecki97/url-shortener/pkg/db/postgresql"
	"github.com/dbielecki97/url-shortener/pkg/db/redis"
	"github.com/dbielecki97/url-shortener/pkg/logger"
	"github.com/dbielecki97/url-shortener/pkg/server"
	"github.com/dbielecki97/url-shortener/pkg/url"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func main() {
	l := logger.New()
	l.Infof("Starting email shortener server on port 8000 ...")

	sanityCheck(l)

	r, rClose := redis.New(l)
	defer rClose()
	p, pClose := postgresql.New(l)
	defer pClose()
	service := url.NewService(r, p)
	s := server.New(l, service)

	srv := &http.Server{
		Handler:      s,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		l.Fatalf("Could not start server: %v", err)
	}
}

func sanityCheck(log *logrus.Logger) {
	keys := []string{
		"REDIS_HOST",
		"POSTGRES_HOST"}

	allPresent := true
	for _, e := range keys {
		ok := checkEnvVariable(e, log)
		if allPresent != false {
			allPresent = ok
		}
	}

	if !allPresent {
		os.Exit(1)
	}
}

func checkEnvVariable(key string, log *logrus.Logger) bool {
	if os.Getenv(key) == "" {
		log.Errorf("Environment variable " + key + " not defined!")
		return false
	}
	return true
}
