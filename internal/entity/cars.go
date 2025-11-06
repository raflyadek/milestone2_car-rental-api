package entity

type Cars struct { 
	Id int `json:"id"`
	Name string `json:"name"`
	PlatNumber string `json:"plat_number"`
	CategoryId int `json:"category_id"`
	Categories Categories `json:"categories_info" gorm:"foreignKey:CategoryId;references:Id"`
	Description string `json:"description"`
	Price float64 `json:"price"`
	Availability bool `json:"availability"`
	AvailabilityUntil string `json:"availability_until"`
}

type Categories struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}