package orderRepo

import (
	"database/sql"
	"fmt"

	"frapuccino/models"
)

func (r orderRepository) GetRepoId(id int) (models.Order, error) {
	var oneOrder models.Order
	query := `
	SELECT
	o.order_id,
	o.customer_name,
	o.status,
	o.created_at,
	oi.product_id,
	oi.quantity
FROM orders o
LEFT JOIN order_items oi ON o.order_id = oi.order_id
WHERE o.order_id = $1;
	`
	rows, err := r.newDB.Db.Query(query, id)
	if err != nil {
		return oneOrder, err
	}
	defer rows.Close()
	oneOrder.Items = []models.OrderItem{}
	for rows.Next() {
		var productId sql.NullInt64
		var quantity, orderId sql.NullInt64
		var customerName, status, createdAt string
		err := rows.Scan(
			&orderId,
			&customerName,
			&status,
			&createdAt,
			&productId,
			&quantity,
		)
		if err != nil {
			return oneOrder, err
		}
		if oneOrder.ID == 0 {
			oneOrder.ID = id
			oneOrder.CustomerName = customerName
			oneOrder.Status = status
			oneOrder.CreatedAt = createdAt
		}

		oneOrder.Items = append(oneOrder.Items, models.OrderItem{
			ProductID: int(productId.Int64),
			Quantity:  int(quantity.Int64),
		})

	}
	if oneOrder.ID == 0 {
		return oneOrder, fmt.Errorf("order with ID %d not found", id)
	}
	return oneOrder, nil
}
