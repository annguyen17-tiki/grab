package store

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"gorm.io/gorm/clause"
)

type ILocationStore interface {
	Save(l *model.Location) error
	Get(l *model.Location) (*model.Location, error)
	Search(conditions *dto.SearchLocations, preloaders ...string) ([]*model.Location, error)
}

type locationStore struct {
	Store
}

func (s locationStore) Save(l *model.Location) error {
	c := clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}},
		UpdateAll: true,
	}

	return s.db.Clauses(c).Model(&model.Location{}).Create(l).Error
}

func (s locationStore) Get(l *model.Location) (*model.Location, error) {
	if err := s.db.Model(&model.Location{}).Where(l).Take(l).Error; err != nil {
		return nil, err
	}

	return l, nil
}

func (s locationStore) Search(conditions *dto.SearchLocations, preloaders ...string) ([]*model.Location, error) {
	query := s.db.Model(&model.Location{})
	if len(conditions.AccountIDs) > 0 {
		query = query.Where("account_id IN (?)", conditions.AccountIDs)
	}

	if len(conditions.GeoHashes) > 0 {
		query = query.Where("geo_hash IN (?)", conditions.GeoHashes)
	}

	if len(conditions.DriverStatuses) > 0 {
		query = query.Where("account_id IN (?)", s.db.Select("account_id").Model(&model.Driver{}).Where("status IN (?)", conditions.DriverStatuses))
	}

	if len(conditions.Vehicles) > 0 {
		query = query.Where("account_id IN (?)", s.db.Select("account_id").Model(&model.Driver{}).Where("vehicle IN (?)", conditions.Vehicles))
	}

	for _, preloader := range preloaders {
		query = query.Preload(preloader)
	}

	var locations []*model.Location
	if err := query.Find(&locations).Error; err != nil {
		return nil, err
	}

	return locations, nil
}
