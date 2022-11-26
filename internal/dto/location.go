package dto

type SaveLocation struct {
	AccountID string  `json:"-"`
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Address   string  `json:"address" validate:"-"`
}

type SearchLocations struct {
	AccountIDs     []string
	GeoHashes      []string
	Vehicles       []string
	DriverStatuses []string
}

type NearestLocations struct {
	Latitude       float64  `form:"latitude" validate:"required"`
	Longitude      float64  `form:"longitude" validate:"required"`
	Vehicles       []string `form:"vehicle" validate:"omitempty,gt=0,unique,dive,oneof=motor car4 car7"`
	DriverStatuses []string `form:"driver_status" validate:"omitempty,gt=0,unique,dive,oneof=free offering book"`
	Radius         float64  `form:"radius" validate:"gte=0"`
}
