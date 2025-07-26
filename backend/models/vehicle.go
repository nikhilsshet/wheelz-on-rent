package models

type Vehicle struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Type         string  `json:"type"` // car or bike
	Model        string  `json:"model"`
	NumberPlate  string  `json:"number_plate"`
	Color        string  `json:"color"`
	Availability bool    `json:"availability"`
	PricePerDay  float64 `json:"price_per_day"`
}
