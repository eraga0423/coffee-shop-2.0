package database

import (
	"time"

	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/models"
)

type SearchFilterRepo interface {
	GetLeftOvers(sortby string, page, pageSize int) (models.LeftOvers, error)
	GetOrderedItems(startDate, endDate time.Time) (map[string]int, error)
	TextSearch(search string, minPrice, maxPrice *float64, filter []string) (models.SearchReports, error)
	OrderedItemsByPeriodDay(month string) (models.PeriodResult, error)
	OrderedItemsByPeriodMonth(year string) (models.PeriodResult, error)
	WriteDBNewOrders(body []models.Order) (*models.Common, error)
}

type searchFilterRepo struct {
	Db *SqlDataBase.DB
}

func NewSearchFilterRepo(db *SqlDataBase.DB) SearchFilterRepo {
	return &searchFilterRepo{Db: db}
}
