package dal

import (
	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/models"
)

// AggregationsRepository defines the interface for reading JSON data for orders and menu items
type AggregationsRepository interface {
	RepositoryTotalSales() (float64, error)
	RepositoryPopularItem() (error, []models.Popular)
}

type aggregationsRepository struct {
	newDB *SqlDataBase.DB
}

// NewAggregationsRepository creates and returns a new instance of aggregationsRepository
func NewAggregationsRepository(db *SqlDataBase.DB) AggregationsRepository {
	return &aggregationsRepository{newDB: db}
}

func (r aggregationsRepository) RepositoryTotalSales() (float64, error) {
	var res float64
	stmt := `
	SELECT 
    COALESCE(SUM(oi.quantity * m.price), 0) AS total_sales
FROM 
    orders o
JOIN 
    order_items oi ON o.order_id = oi.order_id
JOIN 
    menu_items m ON oi.product_id = m.product_id
WHERE 
    o.status = 'close';

`
	row := r.newDB.Db.QueryRow(stmt)
	err := row.Scan(&res)
	if err != nil {
		return 0.0, err
	}

	return res, nil
}

func (r aggregationsRepository) RepositoryPopularItem() (error, []models.Popular) {
	res := []models.Popular{}
	query := `
	SELECT 
		m.name AS popular_item,
		SUM(oi.quantity) AS quantity
	FROM 
		order_items oi
	JOIN 
		menu_items m ON oi.product_id = m.product_id
	GROUP BY 
		m.name
	ORDER BY 
		quantity DESC;
`
	rows, err := r.newDB.Db.Query(query)
	if err != nil {
		return err, nil
	}
	defer rows.Close()
	for rows.Next() {
		var item models.Popular
		err = rows.Scan(&item.PopularSales, &item.Quantity)
		if err != nil {
			return err, nil
		}
		res = append(res, item)
	}
	if err := rows.Err(); err != nil {
		return err, nil
	}
	return nil, res
}
