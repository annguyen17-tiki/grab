package store

import "github.com/annguyen17-tiki/grab/internal/model"

type IDriverStore interface {
	Get(d *model.Driver) (*model.Driver, error)
}

type driverStore struct {
	Store
}

func (s driverStore) Get(d *model.Driver) (*model.Driver, error) {
	if err := s.db.Model(&model.Driver{}).Where(d).Take(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}
