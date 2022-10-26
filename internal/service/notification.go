package service

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
)

func (svc *service) SearchNotifications(input *dto.SearchNotification) ([]*model.Notification, error) {
	return svc.store.Notification().Search(input)
}

func (svc *service) SeenNotification(accountID, notificationID string) error {
	return svc.store.Notification().Seen(accountID, notificationID)
}
