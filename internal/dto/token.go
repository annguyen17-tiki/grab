package dto

type SaveToken struct {
	Token string `json:"token" validate:"required"`
}
