package orderRepo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"frapuccino/models"

	"github.com/lib/pq"
)

func (r *orderRepository) WriteDBNewOrder(body models.Order) error {
	err := r.checkOrdersInMenu(body)
	if err != nil {
		return err
	}

	stmt := `
	INSERT INTO orders (customer_name,  status)
			VALUES ($1, $2)
	RETURNING order_id;
	`
	tx, err := r.newDB.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	row := tx.QueryRow(stmt, body.CustomerName, body.Status)
	err = row.Scan(&body.ID)
	if err != nil {
		return err
	}
	stmt1 := `
	INSERT INTO order_items (order_id, product_id,  quantity)
			VALUES ($1, $2, $3);
	`
	for _, item := range body.Items {
		_, err = tx.Exec(stmt1, body.ID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}
	err = r.CheckIngredients(tx, body)
	if err != nil {
		return err
	}
	return nil
}

// Writes a new order to the JSON file and creates a backup in the reserve copy

func (r *orderRepository) checkOrdersInMenu(body models.Order) error {
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
	rows, err := r.newDB.Db.Query(stmt, pq.Array(productsIds))
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
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) CheckIngredients(tx *sql.Tx, body models.Order) error {
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
		return err
	}
	defer rows.Close()
	var missingIngredients []string
	for rows.Next() {
		var ingrId string
		err := rows.Scan(&ingrId)
		if err != nil {
			return err
		}
		missingIngredients = append(missingIngredients, ingrId)

	}
	if len(missingIngredients) > 0 {
		miss := fmt.Sprintf("Insufficient ingredients: %v", missingIngredients)
		return errors.New(miss)
	}
	err = r.deductInventory(tx, body.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepository) deductInventory(tx *sql.Tx, orderID int) error {
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
		WHERE inventory.ingredient_id = ri.ingredient_id;
	`

	_, err := tx.Exec(stmt, orderID)
	if err != nil {
		return fmt.Errorf("failed to deduct inventory: %w", err)
	}

	return nil
}
