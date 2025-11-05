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
