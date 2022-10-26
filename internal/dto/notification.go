package dto

import "github.com/annguyen17-tiki/grab/internal/model"

type SearchNotification struct {
	AccountID *string `form:"-" validate:"uuid"`
	Status    *string `form:"status" validate:"omitempty,oneof=new seen"`
	*model.Paging
}
