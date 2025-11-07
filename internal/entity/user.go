package entity

type User struct {
	Id               int     `json:"id"`
	Email            string  `json:"email"`
	FullName         string  `json:"full_name"`
	Password         string  `json:"-"`
	// Deposit          float64 `json:"deposit"`
	ValidationCode   string  `json:"-"`
	ValidationStatus bool    `json:"validation_status"`
	Role             string  `json:"-"`
}
