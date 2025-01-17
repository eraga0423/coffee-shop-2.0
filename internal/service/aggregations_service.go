package service

import (
	"frapuccino/internal/dal"
	"frapuccino/models"
)

type AggregationsService interface {
	ServiceTotalSales() (float64, error)
	ServicePopularItems() (error, []models.Popular)
}

type aggregationsService struct {
	aggregationsRepo dal.AggregationsRepository
}

// Initializes and returns a new instance of aggregationsService with the provided repository
func NewAggregationsService(aggregationsRepo dal.AggregationsRepository) AggregationsService {
	return &aggregationsService{aggregationsRepo: aggregationsRepo}
}

// Calculates the total sales from all orders by summing the cost of each order
func (s *aggregationsService) ServiceTotalSales() (float64, error) {
	return s.aggregationsRepo.RepositoryTotalSales()
}

// Finds and returns a sorted list of popular items based on quantities ordered
func (s *aggregationsService) ServicePopularItems() (error, []models.Popular) {
	return s.aggregationsRepo.RepositoryPopularItem()
}
