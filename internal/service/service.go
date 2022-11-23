package service

import (
	"fmt"

	"github.com/annguyen17-tiki/grab/internal/store"
	"github.com/annguyen17-tiki/grab/pkg/notipush"
	"github.com/gomodule/redigo/redis"
)

type service struct {
	store     store.IStore
	notipush  notipush.IService
	redisPool *redis.Pool
	cfg       *Config
}

func New() (IService, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config, err: %v", err)
	}

	s, err := store.New()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare store, err: %v", err)
	}

	n, err := notipush.New()
	if err != nil {
		return nil, fmt.Errorf("failed to init FCM, err: %v", err)
	}

	svc := &service{
		store:    s,
		notipush: n,
		cfg:      cfg,
		redisPool: &redis.Pool{
			MaxActive: 5,
			MaxIdle:   5,
			Wait:      true,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", cfg.RedisURL, redis.DialPassword(cfg.RedisPassword))
			},
		},
	}

	go svc.startJob()
	go svc.startWorkers()

	return svc, nil
}
