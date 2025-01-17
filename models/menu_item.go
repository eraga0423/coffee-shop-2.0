package models

type MenuItem struct {
	ID          int                  `json:"product_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Category    string               `json:"category"`
	Allergens   []string             `json:"allergens"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
}
type MenuItemIngredient struct {
	IngredientID int     `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}
