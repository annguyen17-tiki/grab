package service

import (
	"fmt"

	"github.com/annguyen17-tiki/grab/internal/store"
)

type service struct {
	store store.IStore
	cfg   *Config
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

	svc := &service{
		store: s,
		cfg:   cfg,
	}

	svc.startJob()

	return svc, nil
}
