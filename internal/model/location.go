package model

import "math"

type Location struct {
	AccountID string   `json:"account_id" gorm:"primaryKey"`
	Account   *Account `json:"-"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Address   string   `json:"address"`
	GeoHash   string   `json:"-"`
	Driver    *Driver  `json:"driver" gorm:"foreignKey:AccountID"`
	Trackers
}

// Haversine formula
func (l *Location) DistanceTo(latitude, longitude float64) float64 {
	const earthRadius = 6371
	d := math.Pi / 180
	a := 0.5 - math.Cos((latitude-l.Latitude)*d)/2 + math.Cos(l.Latitude*d)*math.Cos(latitude*d)*(1-math.Cos((longitude-l.Longitude)*d))/2
	return 2 * earthRadius * math.Asin(math.Sqrt(a))
}
