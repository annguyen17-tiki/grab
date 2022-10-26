package model

const (
	RoleAdmin  = "admin"
	RoleDriver = "driver"
	RoleUser   = "user"
)

type Account struct {
	ID        string  `json:"id"`
	Username  string  `json:"-"`
	Password  string  `json:"-"`
	Firstname string  `json:"firstname"`
	Lastname  string  `json:"lastname"`
	Phone     string  `json:"-"`
	Avatar    string  `json:"avatar"`
	Role      string  `json:"role"`
	Driver    *Driver `json:"driver,omitempty"`
	Trackers
}
