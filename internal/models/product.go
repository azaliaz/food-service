package models

type Product struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Mealtype      string  `json:"mealtype"`
	Calories      float64 `json:"calories"`
	Protein       float64 `json:"protein"`
	Fat           float64 `json:"fat"`
	Carbohydrates float64 `json:"carbohydrates"`
	Grams         float64 `json:"grams"`
}
