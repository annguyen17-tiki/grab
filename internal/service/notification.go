package service

import (
	"fmt"

	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/gocraft/work"
)

func (svc *service) SearchNotifications(input *dto.SearchNotification) ([]*model.Notification, error) {
	return svc.store.Notification().Search(input)
}

func (svc *service) SeenNotification(accountID, notificationID string) error {
	return svc.store.Notification().Seen(accountID, notificationID)
}

func (svc *service) notifyNewBookingToDrivers(booking *model.Booking, drivers []*model.Account) error {
	enqueuer := work.NewEnqueuer(model.RedisNamespace, svc.redisPool)

	var notifications []*model.Notification
	for _, driver := range drivers {
		notifications = append(notifications, &model.Notification{
			AccountID: driver.ID,
			Status:    model.NotificationNew,
			Content: map[string]interface{}{
				"title":      "Chuyến xe mới",
				"message":    "Bạn có một yêu cầu đặt xe",
				"booking_id": booking.ID,
			},
		})

		if _, err := enqueuer.Enqueue(model.FCMWorkerTopic, map[string]interface{}{
			"account_id": driver.ID,
			"title":      "Bác tài ơi !!",
			"body":       fmt.Sprintf("%s ơi, có một khách hàng gần bạn", driver.Firstname),
			"link":       fmt.Sprintf("%s/bookings/%s", svc.cfg.WebBaseURL, booking.ID),
		}); err != nil {
			return err
		}
	}

	if err := svc.store.Notification().Create(notifications); err != nil {
		return err
	}

	return nil
}

func (svc *service) notifyConfirmBookingToUser(booking *model.Booking, user, driver *model.Account) error {
	notifications := []*model.Notification{
		{
			AccountID: user.ID,
			Status:    model.NotificationNew,
			Content: map[string]interface{}{
				"title":      "Xác nhận chuyến đi",
				"message":    "Tài xế đã nhận chuyến và đang di chuyển đến chỗ của bạn.",
				"booking_id": booking.ID,
			},
		},
	}

	if err := svc.store.Notification().Create(notifications); err != nil {
		return err
	}

	enqueuer := work.NewEnqueuer(model.RedisNamespace, svc.redisPool)
	if _, err := enqueuer.Enqueue(model.FCMWorkerTopic, map[string]interface{}{
		"account_id": user.ID,
		"title":      "Có xe rồi nha !!",
		"body":       fmt.Sprintf("%s ơi, tài xế %s đang đến đón bạn", user.Firstname, driver.Firstname),
		"link":       fmt.Sprintf("%s/bookings/%s", svc.cfg.WebBaseURL, booking.ID),
	}); err != nil {
		return err
	}

	return nil
}

func (svc *service) notifyTimeoutBookingToUser(booking *model.Booking, user *model.Account) error {
	notifications := []*model.Notification{
		{
			AccountID: user.ID,
			Status:    model.NotificationNew,
			Content: map[string]interface{}{
				"title":      "Chuyến xe bị hủy",
				"message":    "Các tài xế đều đang bận. Vui lòng thử lại sau",
				"booking_id": booking.ID,
			},
		},
	}

	if err := svc.store.Notification().Create(notifications); err != nil {
		return err
	}

	enqueuer := work.NewEnqueuer(model.RedisNamespace, svc.redisPool)
	if _, err := enqueuer.Enqueue(model.FCMWorkerTopic, map[string]interface{}{
		"account_id": user.ID,
		"title":      "Không tìm được tài xế nào !!",
		"body":       fmt.Sprintf("%s ơi, vui lòng thử lại sau nhé", user.Firstname),
		"link":       fmt.Sprintf("%s/bookings/%s", svc.cfg.WebBaseURL, booking.ID),
	}); err != nil {
		return err
	}

	return nil
}

func (svc *service) notifyDoneBookingToUserAndDriver(booking *model.Booking, user, driver *model.Account) error {
	notifications := []*model.Notification{
		{
			AccountID: user.ID,
			Status:    model.NotificationNew,
			Content: map[string]interface{}{
				"title":      "Hoàn thành chuyến đi",
				"message":    "Tài xế đã hoàn thành chuyến đi",
				"booking_id": booking.ID,
			},
		},
		{
			AccountID: user.ID,
			Status:    model.NotificationNew,
			Content: map[string]interface{}{
				"title":      "Hoàn thành chuyến đi",
				"message":    "Bạn đã hoàn thành chuyến đi",
				"booking_id": booking.ID,
			},
		},
	}

	if err := svc.store.Notification().Create(notifications); err != nil {
		return err
	}

	enqueuer := work.NewEnqueuer(model.RedisNamespace, svc.redisPool)
	if _, err := enqueuer.Enqueue(model.FCMWorkerTopic, map[string]interface{}{
		"account_id": user.ID,
		"title":      "Đến nơi rồi !!",
		"body":       fmt.Sprintf("Cám ơn %s đã sử dụng dịch vụ", user.Firstname),
		"link":       fmt.Sprintf("%s/bookings/%s", svc.cfg.WebBaseURL, booking.ID),
	}); err != nil {
		return err
	}

	return nil
}
