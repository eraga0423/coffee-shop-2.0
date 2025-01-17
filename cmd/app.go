package main

import (
	"fmt"
	"log"
	"net/http"

	"frapuccino/internal/dal/SqlDataBase"
	handlefunc "frapuccino/internal/handler/handle_func"
)

// StartServer initializes and starts the HTTP server on the specified port
func StartServer(port string) error {
	mux := http.NewServeMux()
	newdb := SqlDataBase.NewDB()
	if err := newdb.Init(); err != nil {
		fmt.Println(err)
		return err
	}

	handlefunc.FrappuccinoNewHandler(mux, newdb)
	handlefunc.OrderHandler(mux, newdb)
	handlefunc.AggregationHandler(mux, newdb)
	handlefunc.InvHandler(mux, newdb)
	handlefunc.MenuHandler(mux, newdb)

	// Set up server port and log the server start
	port = fmt.Sprintf(":%s", port)
	log.Println("Server started on port:", port)

	// Start the HTTP server
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	return server.ListenAndServe()
}
