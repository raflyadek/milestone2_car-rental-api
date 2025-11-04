package entity

type User struct {
	Id int `json:"id"`
	Email string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Deposit float64 `json:"deposit"`
}

type UserRegister struct {
	Email string `json:"email"`
	FullName string `json:"full_name"`
}