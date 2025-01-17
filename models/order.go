package models

type Order struct {
	ID           int         `json:"order_id"`
	CustomerName string      `json:"customer_name"`
	Items        []OrderItem `json:"items"`
	Status       string      `json:"status"`
	CreatedAt    string      `json:"created_at"`
}

type OrderItem struct {
	ProductID int `json:"menu_item_id"`
	Quantity  int `json:"quantity"`
}

type OrderRequest struct {
	Orders []Order `json:"orders"`
}
