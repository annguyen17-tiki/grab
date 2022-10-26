package model

const (
	OfferNew     = "new"
	OfferAccept  = "accept"
	OfferReject  = "reject"
	OfferConfirm = "confirm"
	OfferTimeout = "timeout"
)

type Offer struct {
	BookingID string `json:"booking_id" gorm:"primaryKey"`
	DriverID  string `json:"driver_id" gorm:"primaryKey"`
	Status    string `json:"status"`
	Trackers
}
