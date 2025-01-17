package database

import (
	"fmt"

	"frapuccino/models"

	"github.com/lib/pq"
)

func (d searchFilterRepo) TextSearch(search string, minPrice, maxPrice *float64, filter []string) (models.SearchReports, error) {
	if len(filter) == 0 || filter == nil || filter[0] == "" {
		filter = []string{"all"}
	}

	var res models.SearchReports
	if contains(filter, "orders") || contains(filter, "all") {
		query := `
		SELECT 
    o.order_id AS id,
    o.customer_name,
    array_agg(mi.name) AS items,
    SUM(mi.price * oi.quantity) AS total,
    ts_rank(to_tsvector(o.customer_name || ' ' || string_agg(mi.name, ' ')), to_tsquery($1)) AS relevance
FROM 
    orders o
JOIN 
    order_items oi ON o.order_id = oi.order_id
JOIN 
    menu_items mi ON oi.product_id = mi.product_id
GROUP BY 
    o.order_id, o.customer_name
HAVING 
    to_tsvector(o.customer_name || ' ' || string_agg(mi.name, ' ')) @@ to_tsquery($1)
ORDER BY 
    relevance DESC;

	`
		rows, err := d.Db.Db.Query(query, search)
		if err != nil {
			return res, err
		}
		defer rows.Close()
		for rows.Next() {
			var order models.ResultOrders
			err := rows.Scan(
				&order.OrdersId,
				&order.CustomerName,
				pq.Array(&order.Items),
				&order.Total,
				&order.Relevance,
			)
			if err != nil {
				return res, err
			}
			res.Orders = append(res.Orders, order)
		}
		err = rows.Err()
		if err != nil {
			fmt.Printf("rows.Err() %v", err)
			return res, err
		}
	}
	if contains(filter, "menu") || contains(filter, "all") {
		query := `
            SELECT 
                m.product_id AS id,
                m.name,
                m.description,
                m.price,
                ts_rank(to_tsvector(m.name || ' ' || m.description), plainto_tsquery($1)) AS relevance
            FROM 
                menu_items m
            WHERE 
                to_tsvector(m.name || ' ' || m.description) @@ plainto_tsquery($1)
                AND ($2::NUMERIC IS NULL OR m.price >= $2)
                AND ($3::NUMERIC IS NULL OR m.price <= $3)
            ORDER BY 
                relevance DESC;
        `
		rows, err := d.Db.Db.Query(query, search, minPrice, maxPrice)
		if err != nil {
			return res, err
		}

		defer rows.Close()
		for rows.Next() {
			var menuItem models.ResultMenu
			err := rows.Scan(
				&menuItem.MenuId,
				&menuItem.Name,
				&menuItem.Description,
				&menuItem.Price,
				&menuItem.Relevance,
			)
			if err != nil {
				return res, err
			}
			res.MenuItem = append(res.MenuItem, menuItem)
		}
	}
	res.TotalMatches = len(res.Orders) + len(res.MenuItem)
	return res, nil
}

func contains(arr []string, item string) bool {
	for _, s := range arr {
		if item == s {
			return true
		}
	}
	return false
}
