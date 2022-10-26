package model

const (
	DriverFree     = "free"
	DriverOffering = "offering"
	DriverBook     = "book"
)

const (
	Motor    = "motor"
	Car4Seat = "car4"
	Car7Seat = "car7"
)

type Driver struct {
	AccountID string `json:"account_id" gorm:"primaryKey"`
	Vehicle   string `json:"vehicle"`
	Status    string `json:"status"`
	Trackers
}
