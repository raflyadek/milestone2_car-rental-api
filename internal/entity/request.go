package entity

type RegisterRequest struct {
	Email string `json:"email" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type EmailValidationRequest struct {
	Code string `json:"validation_code" validate:"required"`
	Email string `json:"email" validate:"required"`
}

type CreateRentalCarsRequest struct {
	Name string `json:"name" validate:"required"`
	PlatNumber string `json:"plat_number" validate:"required"`
	CategoryId int `json:"category_id" validate:"required"`
	Description string `json:"description" validate:"required"`
	PricePerDay float64 `json:"price_per_day" validate:"required"`
	PricePerWeek float64 `json:"price_per_week" validate:"required"`
	PricePerMonth float64 `json:"price_per_month" validate:"required"`
	Availability bool `json:"availability"`	
}