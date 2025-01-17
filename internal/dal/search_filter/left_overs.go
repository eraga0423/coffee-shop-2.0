package database

import (
	"errors"

	"frapuccino/models"
)

func (d searchFilterRepo) GetLeftOvers(sortby string, page, pageSize int) (models.LeftOvers, error) {
	var res models.LeftOvers
	offset := (page - 1) * pageSize
	var total int
	err := d.Db.Db.QueryRow("SELECT COUNT(*) AS total FROM inventory").Scan(&total)
	if err != nil {
		return models.LeftOvers{}, err
	}

	totalPage := (total + pageSize - 1) / pageSize
	if totalPage < page {
		return models.LeftOvers{}, errors.New("page not found")
	}
	if pageSize > total {
		pageSize = total
	}
	query := `
	SELECT 
    i.name AS ingredient_name,
    i.quantity AS ingredient_quantity,
    COALESCE(i.price, 0) AS product_price
FROM 
    inventory i
ORDER BY 
    CASE WHEN $1 = 'quantity' THEN i.quantity END DESC,
    CASE WHEN $1 = 'price' THEN i.price END DESC
LIMIT $2 OFFSET $3;

	`
	var items []models.Data
	rows, err := d.Db.Db.Query(query, sortby, pageSize, offset)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var data models.Data
		err := rows.Scan(
			&data.Name,
			&data.Quantity,
			&data.Price,
		)
		if err != nil {
			return res, err
		}
		items = append(items, data)
	}
	res = models.LeftOvers{
		CurrentPage: page,
		HasNextPage: page < totalPage,
		PageSize:    pageSize,
		Data:        items,
		TotalPages:  totalPage,
	}
	return res, nil
}
