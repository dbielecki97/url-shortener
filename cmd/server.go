package main

import (
	"errors"
	"fmt"
	"github.com/dbielecki97/url-shortener/internal/app"
	"github.com/dbielecki97/url-shortener/internal/store/postgresql"
	"github.com/dbielecki97/url-shortener/internal/store/redis"
	"github.com/dbielecki97/url-shortener/pkg/logger"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

func main() {
	logger.Info("Starting Url Shortener API server ...")
	checkEnvVariables()

	p := postgresql.New()
	r := redis.New()
	s := app.NewServer(app.NewDefaultService(r, p))

	srv := &http.Server{
		Handler:      s,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal("Could not start server: %v", err)
	}
}

func checkEnvVariables() {
	keys := []string{
		"REDIS_HOST",
		"POSTGRES_HOST"}

	allPresent := true
	for _, e := range keys {
		ok := checkEnvVariable(e)
		if allPresent != false {
			allPresent = ok
		}
	}

	if !allPresent {
		logger.Fatal("exiting application", errors.New("configuration not complete"))
		os.Exit(1)
	}
}

func checkEnvVariable(key string) bool {
	if os.Getenv(key) == "" {
		logger.Error("environment variable not set", errors.New(fmt.Sprintf("missing %s env variable", key)))
		return false
	}
	return true
}
