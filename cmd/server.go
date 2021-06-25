package main

import (
	"github.com/dbielecki97/url-shortener/internal/app"
	"github.com/dbielecki97/url-shortener/internal/store/postgresql"
	"github.com/dbielecki97/url-shortener/internal/store/redis"
	"github.com/dbielecki97/url-shortener/pkg/logger"
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

	p, rClose := postgresql.New(l)
	defer rClose()
	r, pClose := redis.New(l)
	defer pClose()
	ser := app.NewDefaultService(r, p)
	s := app.NewServer(l, ser)

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
