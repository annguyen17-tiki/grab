package service

import (
	"sort"
	"time"

	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/mmcloughlin/geohash"
	"gorm.io/gorm"
)

func (svc *service) CreateBooking(input *dto.CreateBooking) (*model.Booking, error) {
	geoHash := geohash.EncodeWithPrecision(input.FromLatitude, input.FromLongitude, svc.cfg.GeoHashPrecision)
	locations, err := svc.store.Location().Search(&dto.SearchLocations{
		GeoHashes:      append(geohash.Neighbors(geoHash), geoHash),
		Vehicles:       input.Vehicles,
		DriverStatuses: []string{model.DriverFree},
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(locations, func(i, j int) bool {
		return locations[i].DistanceTo(input.FromLatitude, input.FromLongitude) < locations[j].DistanceTo(input.FromLatitude, input.FromLongitude)
	})

	var offers []*model.Offer
	var driverIDs []string

	for i, location := range locations {
		if i >= svc.cfg.MaxOffersPerBooking {
			break
		}

		offers = append(offers, &model.Offer{DriverID: location.AccountID, Status: model.OfferNew})
		driverIDs = append(driverIDs, location.AccountID)
	}

	booking := &model.Booking{
		UserID:        input.UserID,
		FromLatitude:  input.FromLatitude,
		FromLongitude: input.ToLongitude,
		FromAddress:   input.FromAddress,
		ToLatitude:    input.ToLatitude,
		ToLongitude:   input.FromLongitude,
		ToAddress:     input.ToAddress,
		Status:        model.BookingNew,
		Offers:        offers,
	}

	if err := svc.store.Booking().Create(booking); err != nil {
		return nil, err
	}

	drivers, err := svc.store.Account().Search(driverIDs)
	if err != nil {
		return nil, err
	}

	if err := svc.notifyNewBookingToDrivers(booking, drivers); err != nil {
		return nil, err
	}

	return booking, nil
}

func (svc *service) AcceptBooking(bookingID, driverID string) error {
	booking, err := svc.store.Booking().Get(&model.Booking{ID: bookingID}, "Offers")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.NewErrBadRequest("invalid booking id: %s", bookingID)
		}
		return err
	}

	if booking.Status != model.BookingNew {
		return model.NewErrBadRequest("booking status is not new")
	}

	offer := booking.OfferForDriver(driverID)
	if offer == nil {
		return model.NewErrForbidden("booking is not offered for this driver")
	}

	driver, err := svc.store.Driver().Get(&model.Driver{AccountID: driverID})
	if err != nil {
		return err
	}

	if driver.Status != model.DriverFree {
		if offer.Status == model.OfferAccept {
			return nil
		}
		return model.NewErrBadRequest("driver is not free")
	}

	if err := svc.store.Booking().Accept(bookingID, driverID); err != nil {
		return err
	}

	if err := svc.confirmBookingIfAny(booking.ID); err != nil {
		return err
	}

	return nil
}

func (svc *service) RejectBooking(bookingID, driverID string) error {
	booking, err := svc.store.Booking().Get(&model.Booking{ID: bookingID}, "Offers")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.NewErrBadRequest("invalid booking id: %s", bookingID)
		}
		return err
	}

	if booking.Status != model.BookingNew {
		return model.NewErrBadRequest("booking status is not new")
	}

	offer := booking.OfferForDriver(driverID)
	if offer == nil {
		return model.NewErrForbidden("booking is not offered for this driver")
	}

	if offer.Status == model.OfferReject {
		return nil
	}

	if err := svc.store.Booking().Reject(bookingID, driverID); err != nil {
		return err
	}

	if err := svc.confirmBookingIfAny(booking.ID); err != nil {
		return err
	}

	return nil
}

func (svc *service) DoneBooking(bookingID, driverID string) error {
	booking, err := svc.store.Booking().Get(&model.Booking{ID: bookingID}, "Offers", "User", "Driver")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.NewErrBadRequest("invalid booking id: %s", bookingID)
		}
		return err
	}

	offer := booking.OfferForDriver(driverID)
	if offer == nil {
		return model.NewErrForbidden("booking is not offered for this driver")
	}

	if offer.Status != model.OfferConfirm {
		return model.NewErrForbidden("booking has not yet confirmed for this driver")
	}

	if booking.Status == model.BookingDone {
		return nil
	}

	if err := svc.store.Booking().Done(bookingID); err != nil {
		return err
	}

	return svc.notifyDoneBookingToUserAndDriver(booking, booking.User, booking.Driver)
}

func (svc *service) GetBooking(id string) (*model.Booking, error) {
	booking, err := svc.store.Booking().Get(&model.Booking{ID: id}, "User", "Driver", "Offers")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewErrNotFound("not found booking: %s", id)
		}
		return nil, err
	}
	return booking, nil
}

func (svc *service) SearchBooking(input *dto.SearchBooking) ([]*model.Booking, error) {
	return svc.store.Booking().Search(input, "User", "Driver")
}

func (svc *service) confirmBookingIfAny(id string) error {
	booking, err := svc.store.Booking().Get(&model.Booking{ID: id}, "Offers", "User")
	if err != nil {
		return err
	}

	if booking.Status != model.BookingNew {
		return nil
	}

	var driverIDs []string
	for _, o := range booking.Offers {
		if o.Status == model.OfferNew {
			return nil
		}

		if o.Status == model.OfferAccept {
			driverIDs = append(driverIDs, o.DriverID)
		}
	}

	if len(driverIDs) == 0 {
		if booking.CreatedAt.Add(svc.cfg.BookingTimeout).After(time.Now().UTC()) {
			return nil
		}

		if err := svc.store.Booking().Timeout(booking.ID); err != nil {
			return err
		}

		return svc.notifyTimeoutBookingToUser(booking, booking.User)
	}

	locations, err := svc.store.Location().Search(&dto.SearchLocations{AccountIDs: driverIDs}, "Account")
	if err != nil {
		return err
	}

	location := booking.NearestDriverLocation(locations)
	if err := svc.store.Booking().Confirm(booking.ID, location.AccountID); err != nil {
		return err
	}

	return svc.notifyConfirmBookingToUser(booking, booking.User, location.Account)
}
