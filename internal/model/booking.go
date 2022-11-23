package model

const (
	BookingNew     = "new"
	BookingConfirm = "confirm"
	BookingDone    = "done"
	BookingTimeout = "timeout"
	BookingCancel  = "cancel"
)

type Booking struct {
	ID            string   `json:"id"`
	UserID        string   `json:"user_id"`
	User          *Account `json:"user"`
	DriverID      *string  `json:"driver_id,omitempty"`
	Driver        *Account `json:"driver"`
	Vehicle       *string  `json:"vehicle,omitempty"`
	FromLatitude  float64  `json:"from_latitude"`
	FromLongitude float64  `json:"from_longitude"`
	ToLatitude    float64  `json:"to_latitude"`
	ToLongitude   float64  `json:"to_longitude"`
	Status        string   `json:"status"`
	Offers        []*Offer `json:"offers"`
	Trackers
}

func (b *Booking) OffersByStatus(status string) []*Offer {
	var offers []*Offer
	for _, o := range b.Offers {
		if o.Status == status {
			offers = append(offers, o)
		}
	}
	return offers
}

func (b *Booking) OfferForDriver(driverID string) *Offer {
	for i, o := range b.Offers {
		if o.DriverID == driverID {
			return b.Offers[i]
		}
	}
	return nil
}

func (b *Booking) NearestDriverLocation(locations []*Location) *Location {
	if len(locations) == 0 {
		return nil
	}

	minDistance := locations[0].DistanceTo(b.FromLatitude, b.FromLongitude)
	nearestLocation := locations[0]

	for i, l := range locations {
		distance := l.DistanceTo(b.FromLatitude, b.FromLongitude)
		if distance < minDistance {
			minDistance = distance
			nearestLocation = locations[i]
		}
	}

	return nearestLocation
}
