package handlefunc

import (
	"net/http"

	"frapuccino/internal/dal"
	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/internal/handler"
	"frapuccino/internal/service"
)

func AggregationHandler(mux *http.ServeMux, newDb SqlDataBase.DB) {
	// Set up Aggregations: repository, service, and handler
	aggregationsRepo := dal.NewAggregationsRepository(&newDb)
	aggregationsService := service.NewAggregationsService(aggregationsRepo)
	aggregationsHandler := handler.NewAggregationsHandler(aggregationsService)
	mux.HandleFunc("GET /reports/total-sales", aggregationsHandler.TotalSales)
	mux.HandleFunc("GET /reports/popular-items", aggregationsHandler.PopularItems)
}
