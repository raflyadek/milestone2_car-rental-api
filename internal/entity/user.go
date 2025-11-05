package entity

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Deposit float64 `json:"deposit"`
	ValidationCode string `json:"-"`
	ValidationStatus bool `json:"validation_status"`
}

type UserInfo struct {
	Id int `json:"id"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
	Deposit float64 `json:"deposit"`
	ValidationStatus bool `json:"validation_status"`
}