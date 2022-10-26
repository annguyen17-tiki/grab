package model

const (
	NotificationNew  = "new"
	NotificationSeen = "seen"
)

type Notification struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`
	Status    string `json:"status"`
	Content   JSON   `json:"content"`
	Trackers
}
