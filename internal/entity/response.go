package entity

type CarsResponse struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	PlatNumber   string     `json:"plat_number"`
	CategoryId   int        `json:"category_id"`
	Categories   Categories `json:"categories_info"`
	Description  string     `json:"description"`
	Price        float64    `json:"price"`
	Availability bool       `json:"availability"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	ValidationStatus bool `json:"validation_status"`
}

type PaymentInfoResponse struct {
	Id     int  `json:"payment_id"`
	UserId int  `json:"user_id"`
	User   User `json:"user_info"`
	CarId  int  `json:"car_id"`
	Car    Cars `json:"cars_info"`
	// RentalPeriod string `json:"rental_period"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	Price      float64 `json:"total_price"`
	Status     bool    `json:"status"`
	ValidUntil string  `json:"valid_until"`
	CreatedAt  string  `json:"created_at"`
}

type PaidPaymentResponse struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
	// User User `json:"user_info"`
	CarId int `json:"car_id"`
	// Car Cars `json:"car_info"`
	TotalDay   int     `json:"total_day"`
	TotalSpent float64 `json:"total_spent"`
	CreatedAt  string  `json:"created_at"`
}

type RentalLogsResponseAdmin struct {
	Id         int     `json:"id"`
	UserId     int     `json:"user_id"`
	CarId      int     `json:"car_id"`
	PaymentId  int     `json:"payment_id"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	TotalDay   int     `json:"total_day"`
	TotalSpent float64 `json:"total_spent"`
	CreatedAt  string  `json:"rental_at"`
}

type RentalLogsResponseUser struct {
	Id         int     `json:"id"`
	UserId     int     `json:"user_id"`
	CarId      int     `json:"car_id"`
	Car Cars `json:"car_info"`
	PaymentId  int     `json:"payment_id"`
	StartDate  string  `json:"start_date"`
	EndDate    string  `json:"end_date"`
	TotalDay   int     `json:"total_day"`
	TotalSpent float64 `json:"total_spent"`
	CreatedAt  string  `json:"rental_at"`
}

type RentalAvailabilityResponse struct {
	// CarId int `json:"car_id"`
	// Car Cars `json:"car_info"`
	Availability bool `json:"availability_book"` 
}