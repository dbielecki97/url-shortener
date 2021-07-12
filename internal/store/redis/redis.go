package redis

import (
	"encoding/json"
	"fmt"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/dbielecki97/url-shortener/pkg/logger"
	"github.com/go-redis/redis"
	"os"
	"strings"
	"time"
)

type Redis struct {
	db *redis.Client
}

func New() *Redis {
	host := os.Getenv("REDIS_HOST")

	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:6379", host),
		Password: "",
		DB:       0,
	})

	if _, err := db.Ping().Result(); err != nil {
		logger.Fatal("Could not ping Redis: %v", err)
	}

	logger.Info("Connected to Redis...")
	return &Redis{db: db}
}

func (r Redis) Save(entry *domain.ShortURL) (*domain.ShortURL, errs.RestErr) {
	bytes, err := json.Marshal(&entry)
	if err != nil {
		logger.Error("Could not marshal ShortURL: %v", err)
		return nil, errs.NewUnexpectedError("unexpected error")
	}

	result := r.db.Set(entry.Code, bytes, time.Minute*2)
	if result.Err() != nil {
		logger.Error("Could not save ShortURL to cache: %v", result.Err())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return entry, nil
}

func (r Redis) Find(code string) (*domain.ShortURL, errs.RestErr) {
	result := r.db.Get(code)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			logger.Info(fmt.Sprintf("Could not find ShortURL with key %v in the cache", code))
			return &domain.ShortURL{}, nil
		}
	}

	jsonString, err := result.Result()
	if err != nil {
		logger.Error("Could not get result: %v", err)
		return nil, errs.NewUnexpectedError("unexpected error while getting result")
	}

	var entry domain.ShortURL
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&entry)
	if err != nil {
		logger.Error("Could not decode ShortURL: %v", err)
		return nil, errs.NewUnexpectedError("unexpected error")
	}

	return &entry, nil
}
