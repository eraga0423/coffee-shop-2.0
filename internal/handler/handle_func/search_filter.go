package handlefunc

import (
	"net/http"

	"frapuccino/internal/dal/SqlDataBase"
	database "frapuccino/internal/dal/search_filter"
	"frapuccino/internal/handler"
	"frapuccino/internal/service"
)

func FrappuccinoNewHandler(mux *http.ServeMux, newdb SqlDataBase.DB) {
	searchRepo := database.NewSearchFilterRepo(&newdb)
	searchService := service.NewSearchFilterHandler(searchRepo)
	searchHandler := handler.NewSearchFilterHandler(searchService)
	mux.HandleFunc("GET /orders/numberOfOrderedItems", searchHandler.NumberOfOrderedItems)
	mux.HandleFunc("GET /reports/search", searchHandler.ReportsSearch)
	mux.HandleFunc("GET /reports/orderedItemsByPeriod", searchHandler.OrderedItemsByPeriodHandle)
	mux.HandleFunc("GET /inventory/getLeftOvers", searchHandler.GetLeftOvers)
	mux.HandleFunc("POST /orders/batch-process", searchHandler.BatchProcessHandler)
}
