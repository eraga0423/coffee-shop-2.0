package orderRepo

import (
	"database/sql"

	"frapuccino/models"
)

func (r orderRepository) ParseOrders() ([]models.Order, error) {
	query := `
	SELECT
		o.order_id,
		o.customer_name,
		o.status,
		o.created_at,
		oi.product_id,
		oi.quantity
	FROM orders o
	LEFT JOIN order_items oi ON o.order_id = oi.order_id;
`

	rows, err := r.newDB.Db.Query(query)
	if err != nil {
		return nil, err
	}
	orderMap := make(map[int]*models.Order)
	allOrders := []models.Order{}
	for rows.Next() {
		orderId := 0
		var customerName, status, createdAt string
		var quantity, productID sql.NullInt64

		err = rows.Scan(
			&orderId,
			&customerName,
			&status,
			&createdAt,
			&productID,
			&quantity,
		)
		if err != nil {
			return nil, err
		}
		if _, exists := orderMap[orderId]; !exists {
			orderMap[orderId] = &models.Order{
				ID:           orderId,
				CustomerName: customerName,
				Status:       status,
				CreatedAt:    createdAt,
				Items:        []models.OrderItem{},
			}
		}

		orderMap[orderId].Items = append(orderMap[orderId].Items, models.OrderItem{
			ProductID: int(productID.Int64),
			Quantity:  int(quantity.Int64),
		})

	}

	for _, order := range orderMap {
		allOrders = append(allOrders, *order)
	}

	return allOrders, nil
}
