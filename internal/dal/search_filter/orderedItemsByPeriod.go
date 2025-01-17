package database

import (
	"fmt"

	"frapuccino/models"
)

func (d searchFilterRepo) OrderedItemsByPeriodDay(month string) (models.PeriodResult, error) {
	query := `
		
            SELECT 
                EXTRACT(DAY FROM o.created_at) AS day,
                COUNT(o.order_id) AS total_orders
            FROM 
                orders o
            WHERE 
                TRIM(TO_CHAR(o.created_at, 'Month')) ILIKE $1
                AND o.status = 'close'
            GROUP BY 
                EXTRACT(DAY FROM o.created_at)
            ORDER BY 
                day;
        `
	rows, err := d.Db.Db.Query(query, month)
	if err != nil {
		return models.PeriodResult{}, err
	}
	defer rows.Close()
	var result models.PeriodResult

	for rows.Next() {
		var day int
		var totalOrders int
		err = rows.Scan(&day, &totalOrders)
		if err != nil {
			return models.PeriodResult{}, err
		}

		result.OrderItems = append(result.OrderItems, map[string]int{fmt.Sprintf("%d", day): totalOrders})
	}
	if err := rows.Err(); err != nil {
		return models.PeriodResult{}, err
	}

	return result, nil
}

func (d searchFilterRepo) OrderedItemsByPeriodMonth(month string) (models.PeriodResult, error) {
	query := `
            SELECT 
                TO_CHAR(o.created_at, 'Month') AS month,
                COUNT(o.order_id) AS total_orders
            FROM 
                orders o
            WHERE 
                EXTRACT(YEAR FROM o.created_at) = $1
                AND o.status = 'close'
            GROUP BY 
                TO_CHAR(o.created_at, 'Month'), EXTRACT(MONTH FROM o.created_at)
            ORDER BY 
                EXTRACT(MONTH FROM o.created_at);
        `
	rows, err := d.Db.Db.Query(query, month)
	if err != nil {
		return models.PeriodResult{}, err
	}
	defer rows.Close()
	var result models.PeriodResult
	for rows.Next() {
		var month string
		var totalOrders int
		err = rows.Scan(&month, &totalOrders)
		if err != nil {
			return models.PeriodResult{}, err
		}
		result.OrderItems = append(result.OrderItems, map[string]int{month: totalOrders})

	}
	if err := rows.Err(); err != nil {
		return models.PeriodResult{}, err
	}

	return result, nil
}
