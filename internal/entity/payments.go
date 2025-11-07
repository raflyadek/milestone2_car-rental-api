package entity

type Payments struct {
	Id int `json:"id"`
	UserId int `json:"user_id"`
	User User `json:"user_info" gorm:"foreignKey:UserId;references:Id"`
	CarId int `json:"car_id"`
	Car Cars `json:"cars_info" gorm:"foreignKey:CarId;references:Id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
	Price float64 `json:"price"`
	Status bool `json:"status_payment"`
	ValidUntil string `json:"valid_until"`
	CreatedAt string `json:"created_at"`
}