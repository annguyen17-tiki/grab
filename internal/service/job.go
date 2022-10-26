package service

import (
	"fmt"
	"time"

	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
)

func (svc *service) startJob() {
	go svc.checkBookingStatus()
}

func (svc *service) checkBookingStatus() {
	for {
		status := model.BookingNew
		bookings, err := svc.store.Booking().Search(&dto.SearchBooking{Status: &status}, "Offers")
		if err != nil {
			fmt.Printf("failed to search bookings, err: %v", err)
		}

		for _, booking := range bookings {
			offersAccepted := booking.OffersByStatus(model.OfferAccept)
			if len(offersAccepted) == 0 {
				if booking.CreatedAt.Add(svc.cfg.BookingTimeout).After(time.Now().UTC()) {
					continue
				}

				if err := svc.store.Booking().Timeout(booking.ID); err != nil {
					fmt.Printf("failed to timeout booking, err: %v", err)
				}

				continue
			}

			var driverIDs []string
			for _, offer := range offersAccepted {
				driverIDs = append(driverIDs, offer.DriverID)
			}

			locations, err := svc.store.Location().Search(&dto.SearchLocations{AccountIDs: driverIDs})
			if err != nil {
				fmt.Printf("failed to search driver locations, err: %v", err)
			}

			location := booking.NearestDriverLocation(locations)
			if err := svc.store.Booking().Confirm(booking.ID, location.AccountID); err != nil {
				fmt.Printf("failed to confirm booking, err: %v", err)
			}
		}

		<-time.After(svc.cfg.CheckBookingTimeoutInterval)
	}
}
