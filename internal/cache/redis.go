package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/response"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

type cache struct {
	client *redis.Client
}

func NewRedis(cfg *config.Cache) *cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
		// TLSConfig: &tls.Config{
		// 	MinVersion: tls.VersionTLS12, // TODO: Uncomment this when we have a valid certificate
		// },
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to cache")

	return &cache{
		client: client,
	}
}

func (c *cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Printf("error during cache set: %v", err)
		return err
	}
	return nil
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, response.ErrCacheNotFound
		}
		log.Printf("error during cache get: %v", err)
		return nil, err
	}
	return []byte(res), nil
}

func (c *cache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("error during cache delete: %v", err)
		return err
	}
	return nil
}

func (c *cache) Close() error {
	err := c.client.Close()
	if err != nil {
		log.Printf("error closing the cache: %v", err)
		return err
	}
	log.Printf("Cache connection closed")
	return nil
}

func GenerateKey(prefix string, requiredParam any, opts ...any) string {
	b := strings.Builder{}
	b.WriteString(prefix)
	b.WriteString(fmt.Sprintf(":%v", requiredParam))
	for _, opt := range opts {
		b.WriteString(fmt.Sprintf(":%v", opt))
	}

	return b.String()
}
