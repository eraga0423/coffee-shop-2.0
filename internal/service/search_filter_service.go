package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	database "frapuccino/internal/dal/search_filter"
	"frapuccino/models"
)

type searchFilterService struct {
	orderService
	searchFilterService database.SearchFilterRepo
}

type SearchFilterService interface {
	NumberOfOrderedItemsService(startDate, endDate string) (map[string]int, error)
	ReportsSearchService(req, filter, minPrice, maxPrice string) (models.SearchReports, error)
	OrderedItemsByPeriodService(Period, month, year string) (models.PeriodResult, error)
	GetLeftOversService(sortby, page, pageSize string) (models.LeftOvers, error)
	BulkOrderProcessingService(orders []models.Order) (*models.Common, error)
}

func NewSearchFilterHandler(searchFilterservice database.SearchFilterRepo) SearchFilterService {
	return &searchFilterService{searchFilterService: searchFilterservice}
}

func (s searchFilterService) NumberOfOrderedItemsService(startDate, endDate string) (map[string]int, error) {
	var timeStartDate, timeEndDate time.Time
	layout := "02.01.2006"
	if startDate != "" {
		parseStartDate, err := time.Parse(layout, startDate)
		if err != nil {
			fmt.Println("error parsing start date")
			return nil, err
		}
		timeStartDate = parseStartDate
	} else if startDate == "" {
		timeStartDate = time.Now().AddDate(-5, 0, 0)
	}
	if endDate != "" {
		parseEndDate, err := time.Parse(layout, endDate)
		if err != nil {
			fmt.Println("error parsing end date")
			return nil, err
		}
		timeEndDate = parseEndDate
	} else if endDate == "" {
		timeEndDate = time.Now()
	}
	return s.searchFilterService.GetOrderedItems(timeStartDate, timeEndDate)
}

func (s searchFilterService) ReportsSearchService(search, filter, minPrice, maxPrice string) (models.SearchReports, error) {
	var result models.SearchReports
	if search == "" {
		return result, errors.New("enter a word")
	}
	arrFilter := strings.Split(filter, ",")
	minPriceFloat := parsePrice(minPrice)

	maxPriceFloat := parsePrice(maxPrice)
	return s.searchFilterService.TextSearch(search, minPriceFloat, maxPriceFloat, arrFilter)
}

func parsePrice(priceStr string) *float64 {
	if priceStr == "" {
		return nil
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return nil
	}
	return &price
}

func (s searchFilterService) OrderedItemsByPeriodService(period, month, year string) (models.PeriodResult, error) {
	var result models.PeriodResult
	result.Period = period
	if period == "day" {
		if month == "" {
			month = time.Now().Format("January")
		}
		res, err := s.searchFilterService.OrderedItemsByPeriodDay(month)
		if err != nil {
			return result, err
		}
		result = res

	} else if period == "month" {
		if year == "" {
			year = time.Now().Format("2006")
		}
		res, err := s.searchFilterService.OrderedItemsByPeriodMonth(year)
		if err != nil {
			return res, err
		}
		result = res
	} else {
		return result, errors.New("please enter a valid period")
	}
	result.Period = period
	result.Year = &year
	result.Month = &month
	return result, nil
}

func (s searchFilterService) GetLeftOversService(sortby, page, pageSize string) (models.LeftOvers, error) {
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		pageInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt <= 0 {
		pageSizeInt = 10
	}

	if sortby != "quantity" && sortby != "price" {
		sortby = "quantity"
	}
	res, err := s.searchFilterService.GetLeftOvers(sortby, pageInt, pageSizeInt)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s searchFilterService) BulkOrderProcessingService(orders []models.Order) (*models.Common, error) {
	for _, order := range orders {
		err := s.orderService.CheckBodyOrder(order)
		if err != nil {
			return nil, err
		}
	}
	return s.searchFilterService.WriteDBNewOrders(orders)
}
