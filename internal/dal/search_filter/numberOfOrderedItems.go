package database

import (
	"errors"
	"time"
)

func (d searchFilterRepo) GetOrderedItems(startDate, endDate time.Time) (map[string]int, error) {
	query := `
	SELECT 
            m.name AS product_name,
            SUM(oi.quantity) AS total_quantity
        FROM 
            order_items oi
        JOIN 
            menu_items m ON oi.product_id = m.product_id
        JOIN 
            orders o ON oi.order_id = o.order_id
        WHERE 
            o.status = 'close'
            AND ($1::timestamp IS NULL OR o.created_at >= $1::timestamp)
			AND ($2::timestamp IS NULL OR o.created_at <= $2::timestamp)

        GROUP BY 
            m.name
        ORDER BY 
            total_quantity DESC;
	`
	rows, err := d.Db.Db.Query(query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make(map[string]int)
	for rows.Next() {
		var productName string
		var totalQuantity int
		err := rows.Scan(
			&productName,
			&totalQuantity,
		)
		if err != nil {
			return nil, err
		}
		result[productName] = totalQuantity
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("no data found")
	}

	return result, nil
}
