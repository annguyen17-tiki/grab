package dto

type CreateAccount struct {
	Username  string `json:"username" validate:"min=6"`
	Password  string `json:"password" validate:"min=6"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Phone     string `json:"phone" validate:"omitempty,phone"`
	Avatar    string `json:"avatar" validate:"omitempty,url"`
	Role      string `json:"role" validate:"oneof=user driver admin"`
}

type UpdateAccount struct {
	ID        string `json:"-" validate:"uuid"`
	Firstname string `json:"firstname" validate:"omitempty"`
	Lastname  string `json:"lastname" validate:"omitempty"`
	Phone     string `json:"phone" validate:"omitempty,phone"`
	Avatar    string `json:"avatar" validate:"omitempty,url"`
}

type LoginInput struct {
	Username string `json:"username" validate:"min=6"`
	Password string `json:"password" validate:"min=6"`
}
