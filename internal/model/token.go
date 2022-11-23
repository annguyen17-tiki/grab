package model

type Token struct {
	AccountID string `gorm:"primaryKey"`
	Token     string
	Trackers
}
