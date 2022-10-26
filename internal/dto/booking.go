package dto

import "github.com/annguyen17-tiki/grab/internal/model"

type CreateBooking struct {
	UserID        string   `json:"-"`
	Vehicles      []string `json:"vehicle" validate:"omitempty,gt=0,unique,dive,oneof=motor car4 car7"`
	FromLatitude  float64  `json:"from_latitude" validate:"required"`
	FromLongitude float64  `json:"from_longitude" validate:"required"`
	ToLatitude    float64  `json:"to_latitude" validate:"required"`
	ToLongitude   float64  `json:"to_longitude" validate:"required"`
}

type SearchBooking struct {
	AccountID *string `form:"account_id" validate:"omitempty,uuid"`
	Status    *string `form:"status" validate:"omitempty,oneof=new confirm done timeout cancel"`
	*model.Paging
}
