package entity

type RentalLogs struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	// User User `json:"user_info" gorm:"foreignKey:UserId;references:id"`
	CarId int `json:"car_id"`
	Car Cars `json:"car_info" gorm:"foreignKey:CarId;references:id"`
	PaymentId int `json:"payment_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	TotalDay int `json:"total_day"`
	TotalSpent float64 `json:"total_spent"`
	CreatedAt string `json:"rental_at"`
}