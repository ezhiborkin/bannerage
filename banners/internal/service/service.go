package service

import (
	"banners/internal/storage/redisC"
	"log/slog"
)

type Service struct {
	log           *slog.Logger
	bannerStorage BannerStorage
	userStorage   UserStorage
	c             *redisC.Cache
}

func New(log *slog.Logger, bannerStorage BannerStorage, userStorage UserStorage, c *redisC.Cache) (*Service, error) {
	return &Service{log: log, bannerStorage: bannerStorage, userStorage: userStorage, c: c}, nil
}
