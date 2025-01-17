package orderRepo

import (
	"database/sql"

	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/models"
)

type OrderRepository interface {
	WriteDBNewOrder(body models.Order) error
	DeleteOldOrder(tx *sql.Tx, id int) error
	UpdateOrder(id int, body models.Order) error
	OrderClose(id int) error
	ParseOrders() ([]models.Order, error)
	DeleteOrder(id int) error
	GetRepoId(id int) (models.Order, error)
	CheckIngredients(tx *sql.Tx, body models.Order) error
}

type orderRepository struct {
	newDB *SqlDataBase.DB
}

// Creates and returns a new instance of jsonOrderRepository
func NewJSONOrderRepository(db *SqlDataBase.DB) OrderRepository {
	return &orderRepository{newDB: db}
}
