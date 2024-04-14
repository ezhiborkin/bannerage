package redisC

import (
	"banners/domain/models"
	"banners/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	*redis.Client
}

func New(port string) (*Cache, error) {
	const op = "storage.redisC.New"

	r := redis.NewClient(&redis.Options{
		Addr:     "redis:" + port,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Cache{r}, nil
}

func (c *Cache) Get(ctx context.Context, key string) (*models.Banner, error) {
	const op = "storage.redisC.Get"

	var banner models.Banner
	bannerJSON, err := c.Client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFoundInCache)
	}

	err = json.Unmarshal(bannerJSON, &banner)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &banner, nil
}

func (c *Cache) Set(ctx context.Context, key string, value models.Banner) error {
	const op = "storage.redisC.Set"

	bannerJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = c.Client.Set(ctx, key, bannerJSON, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
