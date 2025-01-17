package handler

import (
	"encoding/json"
	"net/http"

	"frapuccino/internal/service"
	"frapuccino/models"
)

type AggregationsHandler interface {
	PopularItems(w http.ResponseWriter, r *http.Request)
	TotalSales(w http.ResponseWriter, r *http.Request)
}

type aggregationsHandler struct {
	aggregationsService service.AggregationsService
}

func NewAggregationsHandler(aggregationsService service.AggregationsService) AggregationsHandler {
	return &aggregationsHandler{aggregationsService: aggregationsService}
}

// Handles the HTTP request to retrieve and return total sales data as JSON
func (h *aggregationsHandler) TotalSales(w http.ResponseWriter, r *http.Request) {
	total, err := h.aggregationsService.ServiceTotalSales()
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	ReturnedTotal := models.Total{
		TotalSales: total,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ReturnedTotal)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}

// Handles the HTTP request to retrieve and return popular menu items as JSON
func (h *aggregationsHandler) PopularItems(w http.ResponseWriter, r *http.Request) {
	err, res := h.aggregationsService.ServicePopularItems()
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(res)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}
