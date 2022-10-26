package store

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IBookingStore interface {
	Create(b *model.Booking) error
	Get(b *model.Booking, preloaders ...string) (*model.Booking, error)
	Search(conditions *dto.SearchBooking, preloaders ...string) ([]*model.Booking, error)
	Accept(bookingID, driverID string) error
	Reject(bookingID, driverID string) error
	Timeout(id string) error
	Confirm(bookingID, driverID string) error
	Done(bookingID string) error
}

type bookingStore struct {
	Store
}

func (s bookingStore) Create(b *model.Booking) error {
	return s.db.Model(&model.Booking{}).Create(b).Error
}

func (s bookingStore) Get(b *model.Booking, preloaders ...string) (*model.Booking, error) {
	query := s.db.Model(&model.Booking{})
	for _, preloader := range preloaders {
		query.Preload(preloader)
	}

	if err := query.Where(b).Take(b).Error; err != nil {
		return nil, err
	}

	return b, nil
}

func (s bookingStore) Search(conditions *dto.SearchBooking, preloaders ...string) ([]*model.Booking, error) {
	query := s.db.Model(&model.Booking{})

	if conditions.AccountID != nil {
		query = query.Where(s.db.Where(&model.Booking{UserID: *conditions.AccountID}).Or(&model.Booking{DriverID: conditions.AccountID}))
	}

	if conditions.Status != nil {
		query = query.Where(&model.Booking{Status: *conditions.Status})
	}

	for _, preloader := range preloaders {
		query = query.Preload(preloader)
	}

	if conditions.Paging != nil {
		query = applyPaging(query, conditions.Paging)
	}

	var bookings []*model.Booking
	if err := query.Find(&bookings).Error; err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s bookingStore) Accept(bookingID, driverID string) error {
	txn := func(tx *gorm.DB) error {
		if err := tx.Model(&model.Offer{}).
			Where(&model.Offer{BookingID: bookingID, Status: model.OfferNew}).
			Where(&model.Offer{DriverID: driverID}).
			Updates(&model.Offer{Status: model.OfferAccept}).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Driver{}).
			Where(&model.Driver{AccountID: driverID, Status: model.DriverFree}).
			Updates(&model.Driver{Status: model.DriverOffering}).Error; err != nil {
			return err
		}

		return nil
	}

	return s.db.Transaction(txn)
}

func (s bookingStore) Reject(bookingID, driverID string) error {
	return s.db.Model(&model.Offer{}).
		Where(&model.Offer{BookingID: bookingID, Status: model.OfferNew}).
		Where(&model.Offer{DriverID: driverID}).
		Updates(&model.Offer{Status: model.OfferReject}).
		Error
}

func (s bookingStore) Timeout(id string) error {
	txn := func(tx *gorm.DB) error {
		if err := tx.Model(&model.Booking{}).
			Where(&model.Booking{ID: id, Status: model.BookingNew}).
			Updates(&model.Booking{Status: model.BookingTimeout}).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Offer{}).
			Where(&model.Offer{BookingID: id, Status: model.OfferNew}).
			Updates(&model.Offer{Status: model.OfferTimeout}).
			Error; err != nil {
			return err
		}

		return nil
	}

	return s.db.Transaction(txn)
}

func (s bookingStore) Confirm(bookingID, driverID string) error {
	var driver model.Driver
	if err := s.db.Model(&model.Driver{}).Where(&model.Driver{AccountID: driverID}).Take(&driver).Error; err != nil {
		return err
	}

	txn := func(tx *gorm.DB) error {
		if err := tx.Model(&model.Booking{}).
			Where(&model.Booking{ID: bookingID, Status: model.BookingNew}).
			Updates(&model.Booking{Status: model.BookingConfirm, DriverID: &driverID, Vehicle: &driver.Vehicle}).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Offer{}).
			Where(&model.Offer{BookingID: bookingID, DriverID: driverID, Status: model.OfferAccept}).
			Updates(&model.Offer{Status: model.OfferConfirm}).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Offer{}).
			Where(&model.Offer{BookingID: bookingID, Status: model.OfferNew}).
			Updates(&model.Offer{Status: model.OfferTimeout}).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Driver{}).
			Where(&model.Driver{AccountID: driverID, Status: model.DriverOffering}).
			Updates(&model.Driver{Status: model.DriverBook}).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Driver{}).
			Where(&model.Driver{Status: model.DriverOffering}).
			Where(
				"account_id IN (?)",
				tx.Select("driver_id").Model(&model.Offer{}).Where(&model.Offer{BookingID: bookingID, Status: model.OfferAccept}),
			).
			Updates(&model.Driver{Status: model.DriverFree}).
			Error; err != nil {
			return err
		}

		return nil
	}

	return s.db.Transaction(txn)
}

func (s bookingStore) Done(bookingID string) error {
	txn := func(tx *gorm.DB) error {
		var booking model.Booking
		if err := tx.Clauses(clause.Returning{}).Model(&booking).
			Where(&model.Booking{ID: bookingID, Status: model.BookingConfirm}).
			Updates(&model.Booking{Status: model.BookingDone}).
			Error; err != nil {
			return err
		}

		if booking.DriverID == nil {
			return nil
		}

		if err := tx.Model(&model.Driver{}).
			Where(&model.Driver{AccountID: *booking.DriverID, Status: model.DriverBook}).
			Updates(&model.Driver{Status: model.DriverFree}).Error; err != nil {
			return err
		}

		return nil
	}

	return s.db.Transaction(txn)
}
