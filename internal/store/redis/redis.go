package redis

import (
	"encoding/json"
	"fmt"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
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

func (r Redis) Save(entry *domain.ShortURL) (*domain.ShortURL, error) {
	bytes, err := json.Marshal(&entry)
	if err != nil {
		return nil, errors.Errorf("unexpected error: %v", err)
	}

	result := r.db.Set(entry.Code, bytes, time.Minute*2)
	if result.Err() != nil {
		return nil, errors.Errorf("unexpected database error: %v", result.Err())
	}

	return entry, nil
}

func (r Redis) Find(code string) (*domain.ShortURL, error) {
	result := r.db.Get(code)
	if result.Err() != nil {
		if result.Err() == redis.Nil {
			return nil, domain.NotFoundError{Err: errors.New("could not find entity in the cache")}
		}
	}

	jsonString, err := result.Result()
	if err != nil {
		return nil, errors.Errorf("unexpected database error: %v", err)
	}

	var entry domain.ShortURL
	err = json.NewDecoder(strings.NewReader(jsonString)).Decode(&entry)
	if err != nil {
		return nil, errors.Errorf("unexpected error: %v", err)
	}

	return &entry, nil
}
