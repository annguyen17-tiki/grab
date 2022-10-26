package store

import (
	"github.com/annguyen17-tiki/grab/internal/model"
	"gorm.io/gorm"
)

func applyPaging(tx *gorm.DB, p *model.Paging) *gorm.DB {
	return tx.Limit(p.Limit).Offset(p.Offset).Order(p.Order())
}
