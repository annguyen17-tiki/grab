package store

import (
	"github.com/annguyen17-tiki/grab/internal/model"
	"gorm.io/gorm/clause"
)

type ITokenStore interface {
	Save(token *model.Token) error
	Get(accountID string) (string, error)
}

type tokenStore struct {
	Store
}

func (s tokenStore) Save(token *model.Token) error {
	c := clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}},
		UpdateAll: true,
	}

	return s.db.Model(&model.Token{}).Clauses(c).Create(token).Error
}

func (s tokenStore) Get(accountID string) (string, error) {
	var token model.Token
	if err := s.db.Model(&model.Token{}).Where(&model.Token{AccountID: accountID}).Take(&token).Error; err != nil {
		return "", err
	}

	return token.Token, nil
}
