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
	Price float64 `json:"price" validate:"required"`
	Availability bool `json:"availability"`	
}

type SendEmailValidationRequest struct {
	Email    string `json:"email" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Subject  string `json:"subject" validate:"required"`
	TextPart string `json:"text_part" validate:"required"`
	HtmlPart string `json:"html_part" validate:"required"`
}

type CreatePaymentRequest struct {
	UserId int `json:"user_id" validate:"required"`
	CarId int `json:"car_id" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate string `json:"end_date" validate:"required"`
}