package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type store struct {
	rdb *redis.Client
}

type Store interface {
	SaveShortenedURL(ctx context.Context, _url string) (string, error)
	GetFullURL(ctx context.Context, code string) (string, error)
}

func NewStore(rdb *redis.Client) Store {
	return store{rdb}
}

const shortenedHashMapName = "shortened"

func makeGetCodeError(err error) error {
	return fmt.Errorf(
		"failed to get code from %s hashmap: %w",
		shortenedHashMapName, err,
	)
}

func (s store) SaveShortenedURL(ctx context.Context, _url string) (string, error) {
	var code string
	for range 5 {
		code = genCode()
		if err := s.rdb.HGet(ctx, shortenedHashMapName, code).Err(); err != nil {
			if errors.Is(err, redis.Nil) {
				break
			}
			return "", makeGetCodeError(err)
		}
	}

	if err := s.rdb.HSet(ctx, shortenedHashMapName, code, _url).Err(); err != nil {
		return "", fmt.Errorf(
			"failed to set url %s with code %s in %s hashmap: %w",
			_url, code, shortenedHashMapName, err,
		)
	}
	return code, nil
}

func (s store) GetFullURL(ctx context.Context, code string) (string, error) {
	fullURL, err := s.rdb.HGet(ctx, shortenedHashMapName, code).Result()
	if err != nil {
		return "", makeGetCodeError(err)
	}
	return fullURL, nil
}
