package models

type ProcessedOrder struct {
	OrderId      int      `json:"order_id"`
	CustomerName string   `json:"customer_name"`
	Status       string   `json:"status"`
	Total        *float64 `json:"total"`
	Reason       *string  `json:"reason"`
}
type Summary struct {
	TotalOrders      int               `json:"total_orders"`
	Accepted         int               `json:"accepted"`
	Rejected         int               `json:"rejected"`
	TotalRevenue     float64           `json:"total_revenue"`
	InventoryUpdates []InventoryUpdate `json:"inventory_updates"`
}

type InventoryUpdate struct {
	IngredientId int    `json:"ingredient_id"`
	Name         string `json:"name"`
	QuantityUsed int    `json:"quantity_used"`
	Remaining    int    `json:"remaining"`
}
type Common struct {
	ProccesOrders []ProcessedOrder `json:"process_orders"`
	Summarys      Summary          `json:"summary"`
}
