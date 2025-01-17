package models

type Total struct {
	TotalSales float64 `json:"total_sales"`
}
type Popular struct {
	PopularSales string `json:"popular_item"`
	Quantity     int    `json:"quantity"`
}
