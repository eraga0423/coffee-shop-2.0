package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"frapuccino/internal/service"
	"frapuccino/models"
)

type OrderHandler interface {
	PostOrders(w http.ResponseWriter, r *http.Request)
	GetOrders(w http.ResponseWriter, r *http.Request)
	GetOrdersID(w http.ResponseWriter, r *http.Request)
	PutOrdersID(w http.ResponseWriter, r *http.Request)
	DeleteOrdersID(w http.ResponseWriter, r *http.Request)
	PostOrdersIDClose(w http.ResponseWriter, r *http.Request)
}
type orderHandler struct {
	orderService service.OrderService
}

// Initializes and returns a new instance of orderHandler with the provided service
func NewOrderHandler(orderService service.OrderService) OrderHandler {
	return &orderHandler{orderService: orderService}
}

// Handles the HTTP request to create a new order, validating input and returning success or error
func (h orderHandler) PostOrders(w http.ResponseWriter, r *http.Request) {
	if err := CheckContentType(r); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	body := models.Order{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	err = h.orderService.ServicePostOrders(body)
	if err != nil {

		SendError(w, http.StatusBadRequest, err)
		return
	}

	SendSucces(w, http.StatusCreated, "Order opened")
}

// Handles the HTTP request to retrieve all orders and returns them as JSON
func (h orderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderService.GetOrdersService()
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}

// Handles the HTTP request to retrieve a specific order by ID and returns it as JSON
func (h orderHandler) GetOrdersID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	order, err := h.orderService.GetIDOrdersService(id)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}

// Handles the HTTP request to update a specific order by ID, validating input and updating the order
func (h orderHandler) PutOrdersID(w http.ResponseWriter, r *http.Request) {
	if err := CheckContentType(r); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	body := models.Order{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.orderService.ServicePutOrderID(id, body)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	SendSucces(w, http.StatusOK, "Order updated")
}

// Handles the HTTP request to delete a specific order by ID
func (h orderHandler) DeleteOrdersID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.orderService.ServiceDeleteOrdersID(id); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	SendSucces(w, http.StatusOK, "Order deleted")
}

// Handles the HTTP request to close a specific order by ID
func (h orderHandler) PostOrdersIDClose(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	if err := h.orderService.CloseOrder(id); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	SendSucces(w, http.StatusOK, "Order closed")
}
