package handlefunc

import (
	"net/http"

	"frapuccino/internal/dal"
	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/internal/handler"
	"frapuccino/internal/service"
)

func MenuHandler(mux *http.ServeMux, newDb SqlDataBase.DB) {
	// Set up Menu: repository, service, and handler
	menuRepo := dal.NewJSONMenuRepository(&newDb)
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)
	mux.HandleFunc("POST /menu", menuHandler.PostMenu)
	mux.HandleFunc("GET /menu", menuHandler.GetMenu)
	mux.HandleFunc("GET /menu/{id}", menuHandler.GetMenuID)
	mux.HandleFunc("PUT /menu/{id}", menuHandler.PutMenuID)
	mux.HandleFunc("DELETE /menu/{id}", menuHandler.DeleteMenuID)
}
