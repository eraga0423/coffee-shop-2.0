package orderRepo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"frapuccino/models"
)

func (r *orderRepository) UpdateOrder(id int, body models.Order) error {
	tx, err := r.newDB.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			log.Printf("Transaction rollback due to error: %v", err)
			tx.Rollback()

		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if err = r.CheckStatus(id); err != nil {
		return fmt.Errorf("cannot update closed order: %w", err)
	}

	orderUpdateQuery := `UPDATE orders 
							 SET customer_name = $1, status = $2 
							 WHERE order_id = $3`
	_, err = tx.Exec(orderUpdateQuery, body.CustomerName, body.Status, id)
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	deleteItemsQuery := `DELETE FROM order_items WHERE order_id = $1 RETURNING product_id, quantity`
	rows, err := tx.Query(deleteItemsQuery, id)
	var plusOrders []models.OrderItem
	for rows.Next() {
		var quantity, productId int
		rows.Scan(
			&productId,
			&quantity,
		)
		plusOrders = append(plusOrders, models.OrderItem{
			ProductID: productId,
			Quantity:  quantity,
		})
	}

	if err != nil {
		return fmt.Errorf("failed to delete old order items: %w", err)
	}
	err = r.PlusInventory(plusOrders, tx)
	if err != nil {
		return err
	}
	orderItemsQuery := `INSERT INTO order_items (order_id, product_id, quantity) 
							VALUES ($1, $2, $3)`
	for _, item := range body.Items {
		_, err := tx.Exec(orderItemsQuery, id, item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}
	body.ID = id
	err = r.CheckIngredients(tx, body)
	if err != nil {
		return fmt.Errorf("failed to check ingredients: %w", err)
	}
	return nil
}

func (r *orderRepository) DeleteOldOrder(tx *sql.Tx, id int) error {
	query := `
DELETE FROM orders
WHERE order_id =$1;
	`
	exists, err := tx.Exec(query, id)
	if err != nil {
		return err
	}
	oneRow, err := exists.RowsAffected()
	if err != nil {
		return err
	}
	if oneRow == 0 {
		return errors.New("id is incorect")
	}

	return nil
}

func (r *orderRepository) RestockInventory(tx *sql.Tx, id int) error {
	restockStmt := `
	WITH used_ingredients AS (
		SELECT
			mi.ingredient_id,
			SUM(mi.quantity * oi.quantity) AS used_quantity
		FROM menu_item_ingredients mi
		JOIN order_items oi ON mi.product_id = oi.product_id
		WHERE oi.order_id = $1
		GROUP BY mi.ingredient_id
	)
	UPDATE inventory
	SET quantity = quantity + ui.used_quantity
	FROM used_ingredients ui
	WHERE inventory.ingredient_id = ui.ingredient_id;
	`
	_, err := tx.Exec(restockStmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) CheckStatus(id int) error {
	stmt := `
	SELECT status FROM orders WHERE order_id = $1;
	`
	status := ""
	row := r.newDB.Db.QueryRow(stmt, id)
	err := row.Scan(&status)
	if err != nil {
		return err
	}
	if status != "open" {
		return errors.New("order status does not match")
	}
	return nil
}

func (r *orderRepository) PlusInventory(plusData []models.OrderItem, tx *sql.Tx) error {
	var plusInv []models.MenuItemIngredient

	for _, orderItem := range plusData {
		stmt := `
		SELECT 
		ingredient_id,
		quantity
		FROM menu_item_ingredients
		WHERE product_id = $1
		`
		rows, err := tx.Query(stmt, orderItem.ProductID)
		if err != nil {
			return err
		}

		for rows.Next() {
			var ingId int
			var quantity float64
			err := rows.Scan(
				&ingId,
				&quantity,
			)
			if err != nil {
				return err
			}
			plusInv = append(plusInv, models.MenuItemIngredient{
				IngredientID: ingId,
				Quantity:     quantity * float64(orderItem.Quantity),
			})

		}

	}
	for _, oneInv := range plusInv {
		stmt := `
	UPDATE inventory
	SET quantity = quantity + $1
	WHERE ingredient_id = $2;

	`
		_, err := tx.Exec(stmt, oneInv.Quantity, oneInv.IngredientID)
		if err != nil {
			return err
		}

	}

	return nil
}
