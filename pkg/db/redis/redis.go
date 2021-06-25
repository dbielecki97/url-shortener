package redis

import (
	"encoding/json"
	"fmt"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/dbielecki97/url-shortener/pkg/url"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type Redis struct {
	db  *redis.Client
	log *logrus.Logger
}

func New(log *logrus.Logger) (*Redis, func()) {
	host := os.Getenv("REDIS_HOST")

	db := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:6379", host),
		Password: "",
		DB:       0,
	})

	if _, err := db.Ping().Result(); err != nil {
		log.Fatalf("Could not ping Redis: %v", err)
	}

	closeFn := func() {
		err := db.Close()
		if err != nil {
			log.Printf("Could not close Redis: %v", err)
		}
	}
	log.Println("Connected to Redis...")
	return &Redis{db: db, log: log}, closeFn
}

func (r Redis) Save(entry *url.Entry) (*url.Entry, *errs.AppError) {
	bytes, err := json.Marshal(&entry)
	if err != nil {
		r.log.Errorf("Could not marshal Entry: %v", err)
		return nil, errs.NewUnexpectedError("unexpected error")
	}

	result := r.db.Set(entry.Code, bytes, time.Hour*24)
	if result.Err() != nil {
		r.log.Errorf("Could not save Entry to cache: %v", result.Err())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return entry, nil
}

func (r Redis) Find(code string) (*url.Entry, *errs.AppError) {
	result := r.db.Get(code)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			r.log.Warnf("Could not find Entry with key %v in the cache", code)
			return nil, errs.NewCacheMissError("could not found Entry in cache with provided code")
		}
	}

	jsonString, err := result.Result()
	if err != nil {
		r.log.Infof("Could not get result: %v", err)
		return nil, errs.NewUnexpectedError("unexpected error while getting result")
	}

	var entry url.Entry
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&entry)
	if err != nil {
		r.log.Infof("Could not decode Entry: %v", err)
		return nil, errs.NewUnexpectedError("unexpected error")
	}

	return &entry, nil
}
