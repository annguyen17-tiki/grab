package service

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
)

type IService interface {
	CreateAccount(input *dto.CreateAccount) error
	UpdateAccount(input *dto.UpdateAccount) error
	Login(input *dto.LoginInput) (string, error)
	GetAccount(id string) (*model.Account, error)
	GetAccountByPhone(phone string) (*model.Account, error)

	SaveLocation(input *dto.SaveLocation) error
	GetLocation(accountID string) (*model.Location, error)
	NearestLocations(input *dto.NearestLocations) ([]*model.Location, error)

	CreateBooking(input *dto.CreateBooking) (*model.Booking, error)
	AcceptBooking(bookingID, driverID string) error
	RejectBooking(bookingID, driverID string) error
	DoneBooking(bookingID, driverID string) error
	GetBooking(id string) (*model.Booking, error)
	SearchBooking(input *dto.SearchBooking) ([]*model.Booking, error)

	SearchNotifications(input *dto.SearchNotification) ([]*model.Notification, error)
	SeenNotification(accountID, notificationID string) error

	SaveToken(accountID, token string) error
}
