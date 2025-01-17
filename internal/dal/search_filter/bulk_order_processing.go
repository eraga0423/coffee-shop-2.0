package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"frapuccino/models"

	"github.com/lib/pq"
)

func (r *searchFilterRepo) WriteDBNewOrders(bodies []models.Order) (*models.Common, error) {
	var processOrders []models.ProcessedOrder
	var summary models.Summary
	var inventoryUpdates []models.InventoryUpdate
	var totalRevenue float64
	var accepted, rejected int
	tx, err := r.Db.Db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	for _, body := range bodies {
		err := r.checkOrdersInMenu(body)
		if err != nil {
			stringErr := err.Error()
			processOrders = append(processOrders, models.ProcessedOrder{
				OrderId:      0,
				CustomerName: body.CustomerName,
				Status:       "rejected",
				Total:        nil,
				Reason:       &stringErr,
			})
			rejected++
			continue
		}
		orderID, err := r.insertOrder(tx, body)
		if err != nil {
			stringErr := err.Error()
			processOrders = append(processOrders, models.ProcessedOrder{
				OrderId:      orderID,
				CustomerName: body.CustomerName,
				Status:       "rejected",
				Total:        nil,
				Reason:       &stringErr,
			})
			rejected++
			continue
		}
		body.ID = orderID
		err = r.insertOrderItems(tx, body)
		if err != nil {
			stringErr := err.Error()
			processOrders = append(processOrders, models.ProcessedOrder{
				OrderId:      body.ID,
				CustomerName: body.CustomerName,
				Status:       "rejected",
				Total:        nil,
				Reason:       &stringErr,
			})
			rejected++
			continue
		}
		updates, err := r.checkIngredients(tx, body)
		inventoryUpdates = append(inventoryUpdates, updates...)
		if err != nil {
			stringErr := err.Error()
			processOrders = append(processOrders, models.ProcessedOrder{
				OrderId:      0,
				CustomerName: body.CustomerName,
				Status:       "rejected",
				Total:        nil,
				Reason:       &stringErr,
			})
			rejected++
			continue
		}
		total, err := r.calculateOrderTotal(tx, body.ID)
		if err != nil {
			stringErr := err.Error()
			processOrders = append(processOrders, models.ProcessedOrder{
				OrderId:      body.ID,
				CustomerName: body.CustomerName,
				Status:       "rejected",
				Total:        nil,
				Reason:       &stringErr,
			})
			rejected++
			continue
		}
		totalRevenue += total
		processOrders = append(processOrders, models.ProcessedOrder{
			OrderId:      body.ID,
			CustomerName: body.CustomerName,
			Status:       "accepted",
			Total:        &total,
			Reason:       nil,
		})
		accepted++

	}
	summary = models.Summary{
		TotalOrders:      len(bodies),
		Accepted:         accepted,
		Rejected:         rejected,
		TotalRevenue:     totalRevenue,
		InventoryUpdates: inventoryUpdates,
	}
	return &models.Common{
		ProccesOrders: processOrders,
		Summarys:      summary,
	}, nil
}

// Writes a new order to the JSON file and creates a backup in the reserve copy

func (r *searchFilterRepo) checkOrdersInMenu(body models.Order) error {
	stmt := `
	WITH new_order AS (
		SELECT UNNEST($1::INT[]) AS product_id
	)
	SELECT product_id
	FROM new_order
	WHERE NOT EXISTS (
		SELECT 1
		FROM menu_items
		WHERE menu_items.product_id = new_order.product_id
	);
	`

	productsIds := []int{}
	for _, item := range body.Items {
		productsIds = append(productsIds, item.ProductID)
	}
	rows, err := r.Db.Db.Query(stmt, pq.Array(productsIds))
	if err != nil {
		return err
	}
	defer rows.Close()
	missingProducts := []string{}

	for rows.Next() {
		var missingProduct string
		err := rows.Scan(&missingProduct)
		if err != nil {
			return err
		}
		missingProducts = append(missingProducts, missingProduct)

	}
	if len(missingProducts) > 0 {
		miss := fmt.Sprintf("These items are not in menu %s", strings.Join(missingProducts, ", "))
		return errors.New(miss)
	}
	return nil
}

func (r *searchFilterRepo) checkIngredients(tx *sql.Tx, body models.Order) ([]models.InventoryUpdate, error) {
	stmt := `
	WITH required_ingredients AS (
		SELECT mi.ingredient_id, mi.quantity * oi.quantity AS required_quantity
		FROM menu_item_ingredients mi
		JOIN order_items oi ON mi.product_id = oi.product_id
		WHERE oi.order_id = $1
	),
	insufficient_ingredients AS (
		SELECT ri.ingredient_id, i.quantity AS available_quantity, ri.required_quantity
		FROM required_ingredients ri
		JOIN inventory i ON ri.ingredient_id = i.ingredient_id
		WHERE ri.required_quantity > i.quantity
	)
	SELECT ingredient_id
	FROM insufficient_ingredients;
	`
	rows, err := tx.Query(stmt, body.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var missingIngredients []string
	for rows.Next() {
		var ingrId string
		err := rows.Scan(&ingrId)
		if err != nil {
			return nil, err
		}
		missingIngredients = append(missingIngredients, ingrId)

	}
	if len(missingIngredients) > 0 {
		miss := fmt.Sprintf("Insufficient ingredients: %v", missingIngredients)
		return nil, errors.New(miss)
	}
	updates, err := r.deductInventory(tx, body.ID)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (r *searchFilterRepo) deductInventory(tx *sql.Tx, orderID int) ([]models.InventoryUpdate, error) {
	stmt := `
		UPDATE inventory
		SET quantity = quantity - ri.required_quantity
		FROM (
			SELECT mi.ingredient_id, SUM(mi.quantity * oi.quantity) AS required_quantity
			FROM menu_item_ingredients mi
			JOIN order_items oi ON mi.product_id = oi.product_id
			WHERE oi.order_id = $1
			GROUP BY mi.ingredient_id
		) AS ri
		WHERE inventory.ingredient_id = ri.ingredient_id
		RETURNING inventory.ingredient_id, inventory.name, ri.required_quantity, inventory.quantity;
	`

	rows, err := tx.Query(stmt, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	updates := []models.InventoryUpdate{}
	for rows.Next() {
		var update models.InventoryUpdate
		err := rows.Scan(&update.IngredientId, &update.Name, &update.QuantityUsed, &update.Remaining)
		if err != nil {
			return nil, err
		}
		updates = append(updates, update)

	}
	if len(updates) == 0 {
		return nil, errors.New("No inventory updates")
	}
	return updates, nil
}

func (r *searchFilterRepo) insertOrder(tx *sql.Tx, body models.Order) (int, error) {
	stmt := `
		INSERT INTO orders (customer_name, status)
		VALUES ($1, $2)
		RETURNING order_id;
	`
	var orderID int
	row := tx.QueryRow(stmt, body.CustomerName, "open")
	err := row.Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (r *searchFilterRepo) insertOrderItems(tx *sql.Tx, body models.Order) error {
	stmt := `
		INSERT INTO order_items (order_id, product_id, quantity)
		VALUES ($1, $2, $3);
	`
	for _, item := range body.Items {
		_, err := tx.Exec(stmt, body.ID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}

	}
	return nil
}

func (r *searchFilterRepo) calculateOrderTotal(tx *sql.Tx, orderID int) (float64, error) {
	stmt := `
		SELECT SUM(mi.price * oi.quantity)
		FROM menu_items mi
		JOIN order_items oi ON mi.product_id = oi.product_id
		WHERE oi.order_id = $1;
	`
	var total float64
	row := tx.QueryRow(stmt, orderID)
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
