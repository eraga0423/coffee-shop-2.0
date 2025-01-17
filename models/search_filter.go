package models

type SearchReports struct {
	MenuItem     []ResultMenu   `json:"menu_items"`
	Orders       []ResultOrders `json:"orders"`
	TotalMatches int            `json:"total_matches"`
}
type ResultMenu struct {
	MenuId      string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Relevance   float64 `json:"relevance"`
}
type ResultOrders struct {
	OrdersId     string   `json:"id"`
	CustomerName string   `json:"customer_name"`
	Items        []string `json:"items"`
	Total        float64  `json:"total"`
	Relevance    float64  `json:"relevance"`
}
type PeriodResult struct {
	Period     string           `json:"period"`
	Month      *string          `json:"month"`
	Year       *string          `json:"year"`
	OrderItems []map[string]int `json:"order_items"`
}
type LeftOvers struct {
	CurrentPage int    `json:"current_page"`
	HasNextPage bool   `json:"has_next_page"`
	PageSize    int    `json:"page_size"`
	TotalPages  int    `json:"total_pages"`
	Data        []Data `json:"data"`
}
type Data struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
