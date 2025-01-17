package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"frapuccino/internal/service"
	"frapuccino/models"
)

type SearchFilterHandler interface {
	NumberOfOrderedItems(w http.ResponseWriter, r *http.Request)
	ReportsSearch(w http.ResponseWriter, r *http.Request)
	OrderedItemsByPeriodHandle(w http.ResponseWriter, r *http.Request)
	GetLeftOvers(w http.ResponseWriter, r *http.Request)
	BatchProcessHandler(w http.ResponseWriter, r *http.Request)
}

type searchFilterHandler struct {
	searchFilterService service.SearchFilterService
}

func NewSearchFilterHandler(searchFilterservice service.SearchFilterService) SearchFilterHandler {
	return &searchFilterHandler{searchFilterService: searchFilterservice}
}

func (h *searchFilterHandler) NumberOfOrderedItems(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
	fmt.Println(startDate, endDate)
	results, err := h.searchFilterService.NumberOfOrderedItemsService(startDate, endDate)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	slog.Info("Successfully retrieved number of ordered items")
}

func (h *searchFilterHandler) ReportsSearch(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("q")
	filter := r.URL.Query().Get("filter")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")
	data, err := h.searchFilterService.ReportsSearchService(search, filter, minPrice, maxPrice)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	slog.Info("Successfully retrieved reports")
}

func (h *searchFilterHandler) OrderedItemsByPeriodHandle(w http.ResponseWriter, r *http.Request) {
	dayMonth := r.URL.Query().Get("period")
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")
	if dayMonth == "" {
		SendError(w, http.StatusBadRequest, errors.New("period parameter is required"))
		return
	}
	result, err := h.searchFilterService.OrderedItemsByPeriodService(dayMonth, month, year)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *searchFilterHandler) GetLeftOvers(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")
	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("pageSize")
	res, err := h.searchFilterService.GetLeftOversService(sortBy, page, pageSize)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	slog.Info("Successfully retrieved left overs")
}

func (h *searchFilterHandler) BatchProcessHandler(w http.ResponseWriter, r *http.Request) {
	var orders models.OrderRequest
	err := json.NewDecoder(r.Body).Decode(&orders)
	if err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}
	res, err := h.searchFilterService.BulkOrderProcessingService(orders.Orders)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}
	SendSucces(w, http.StatusCreated, "Successfully processed orders")
}
