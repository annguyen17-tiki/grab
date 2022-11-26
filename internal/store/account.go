package store

import (
	"github.com/annguyen17-tiki/grab/internal/model"
)

type IAccountStore interface {
	Create(a *model.Account) error
	Update(a *model.Account) error
	Get(a *model.Account, preloaders ...string) (*model.Account, error)
	Search(ids []string) ([]*model.Account, error)
}

type accountStore struct {
	Store
}

func (s accountStore) Create(a *model.Account) error {
	return s.db.Model(&model.Account{}).Create(a).Error
}

func (s accountStore) Update(a *model.Account) error {
	return s.db.Model(&model.Account{}).Where(&model.Account{ID: a.ID}).Updates(a).Error
}

func (s accountStore) Get(a *model.Account, preloaders ...string) (*model.Account, error) {
	query := s.db.Model(&model.Account{})
	for _, preloader := range preloaders {
		query = query.Preload(preloader)
	}

	if err := query.Where(a).Take(a).Error; err != nil {
		return nil, err
	}

	return a, nil
}

func (s accountStore) Search(ids []string) ([]*model.Account, error) {
	var accounts []*model.Account
	if err := s.db.Model(&model.Account{}).Where("id IN (?)", ids).Find(&accounts).Error; err != nil {
		return nil, err
	}

	return accounts, nil
}
