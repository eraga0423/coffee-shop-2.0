package handlefunc

import (
	"net/http"

	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/internal/dal/orderRepo"
	"frapuccino/internal/handler"
	"frapuccino/internal/service"
)

func OrderHandler(mux *http.ServeMux, newDb SqlDataBase.DB) {
	// Set up Orders: repository, service, and handler

	orderRepo := orderRepo.NewJSONOrderRepository(&newDb)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)
	mux.HandleFunc("POST /orders", orderHandler.PostOrders)
	mux.HandleFunc("GET /orders", orderHandler.GetOrders)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetOrdersID)
	mux.HandleFunc("PUT /orders/{id}", orderHandler.PutOrdersID)
	mux.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrdersID)
	mux.HandleFunc("POST /orders/{id}/close", orderHandler.PostOrdersIDClose)
}
