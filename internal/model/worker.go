package model

import "encoding/json"

const (
	RedisNamespace = "grab-noti"
	FCMWorkerTopic = "fcm"
)

type WorkerPayload struct {
	AccountID string `json:"account_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Link      string `json:"link"`
}

func CreateWorkerPayload(from map[string]interface{}) (*WorkerPayload, error) {
	b, err := json.Marshal(from)
	if err != nil {
		return nil, err
	}

	var payload WorkerPayload
	err = json.Unmarshal(b, &payload)
	return &payload, err
}
