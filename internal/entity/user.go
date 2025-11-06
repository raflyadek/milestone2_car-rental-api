package entity

type User struct {
	Id               int     `json:"id"`
	Email            string  `json:"email"`
	FullName         string  `json:"full_name"`
	Password         string  `json:"password"`
	Deposit          float64 `json:"deposit"`
	ValidationCode   string  `json:"-"`
	ValidationStatus bool    `json:"validation_status"`
	Role             string  `json:"role"`
}

type UserResponse struct {
	Id               int     `json:"id"`
	Email            string  `json:"email"`
	FullName         string  `json:"full_name"`
	Deposit          float64 `json:"deposit"`
	ValidationStatus bool    `json:"validation_status"`
}

type SendEmailValidationRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Subject  string `json:"subject"`
	TextPart string `json:"text_part"`
	HtmlPart string `json:"html_part"`
}
