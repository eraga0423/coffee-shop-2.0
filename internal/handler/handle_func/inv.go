package handlefunc

import (
	"net/http"

	"frapuccino/internal/dal"
	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/internal/handler"
	"frapuccino/internal/service"
)

func InvHandler(mux *http.ServeMux, newDb SqlDataBase.DB) {
	// Set up Inventory: repository, service, and handler
	invRepo := dal.NewJSONInvRepository(&newDb)
	invService := service.NewInvService(invRepo)
	invHandler := handler.NewInvHandler(invService)
	mux.HandleFunc("POST /inventory", invHandler.PostInv)
	mux.HandleFunc("GET /inventory", invHandler.GetInv)
	mux.HandleFunc("GET /inventory/{id}", invHandler.GetInvID)
	mux.HandleFunc("PUT /inventory/{id}", invHandler.PutInvID)
	mux.HandleFunc("DELETE /inventory/{id}", invHandler.DeleteInvID)
}
