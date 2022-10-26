package store

import (
	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
)

type INotificationStore interface {
	Create(ns []*model.Notification) error
	Search(conditions *dto.SearchNotification) ([]*model.Notification, error)
	Seen(accountID, notificationID string) error
}

type notificationStore struct {
	Store
}

func (s notificationStore) Create(ns []*model.Notification) error {
	return s.db.Model(&model.Notification{}).Create(&ns).Error
}

func (s notificationStore) Search(conditions *dto.SearchNotification) ([]*model.Notification, error) {
	query := s.db.Model(&model.Notification{})
	if conditions.AccountID != nil {
		query = query.Where(&model.Notification{AccountID: *conditions.AccountID})
	}

	if conditions.Status != nil {
		query = query.Where(&model.Notification{Status: *conditions.Status})
	}

	if conditions.Paging != nil {
		query = applyPaging(query, conditions.Paging)
	}

	var notifications []*model.Notification
	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s notificationStore) Seen(accountID, notificationID string) error {
	return s.db.Model(&model.Notification{}).
		Where(&model.Notification{ID: notificationID, AccountID: accountID, Status: model.NotificationNew}).
		Updates(&model.Notification{Status: model.NotificationSeen}).Error
}
