package entity

type CarsResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	PlatNumber string `json:"plat_number"`
	CategoryId int `json:"category_id"`
	Categories Categories `json:"categories_info"`
	Description string `json:"description"`
	PricePerDay float64 `json:"price_per_day"`
	PricePerWeek float64 `json:"price_per_week"`
	PricePerMonth float64 `json:"price_per_month"`
	Availability bool `json:"availability"`
}