package dal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/models"

	"github.com/lib/pq"
)

type (
	MenuRepository interface {
		PostRepoMenu(content models.MenuItem) error
		UpdateMenu(id int, content models.MenuItem) error
		DeleteMenuItem(id int) error
		GetMenuRepo() ([]models.MenuItem, error)
		GetMenuItemID(id int) (models.MenuItem, error)
	}
	jsonMenuRepository struct {
		newDB *SqlDataBase.DB
	}
)

// Creates and returns a new instance of jsonMenuRepository
func NewJSONMenuRepository(db *SqlDataBase.DB) MenuRepository {
	return &jsonMenuRepository{newDB: db}
}

func (m *jsonMenuRepository) PostRepoMenu(content models.MenuItem) error {
	tx, err := m.newDB.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	productId := 0
	stmt := `
	INSERT INTO menu_items (name, description, price, category, allergens)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING product_id;
	`
	row := tx.QueryRow(stmt, content.Name, content.Description, content.Price, content.Category, pq.Array(content.Allergens))
	err = row.Scan(&productId)
	if err != nil {
		return err
	}
	ingredientStmt := `
INSERT INTO menu_item_ingredients (product_id, ingredient_id, quantity)
VALUES ($1, $2, $3)
`
	for _, ingredient := range content.Ingredients {
		_, err := tx.Exec(ingredientStmt, productId, ingredient.IngredientID, ingredient.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *jsonMenuRepository) UpdateMenu(id int, content models.MenuItem) error {
	tx, err := m.newDB.Db.Begin()
	if err != nil {
		return err
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
	deleteMenuQuery := `
UPDATE menu_items
SET name = $1,    description = $2, price = $3, category = $4, allergens = $5
WHERE product_id = $6`
	_, err = tx.Exec(deleteMenuQuery, content.Name, content.Description, content.Price, content.Category, pq.Array(content.Allergens), id)
	if err != nil {
		return err
	}
	deleteIngredientQuery := `
	DELETE FROM menu_item_ingredients WHERE product_id = $1
	`
	_, err = tx.Exec(deleteIngredientQuery, id)
	if err != nil {
		return err
	}
	insertIngredientsQuery := `
	INSERT INTO menu_item_ingredients (product_id, ingredient_id, quantity)
	VALUES ($1, $2, $3)
	`
	for _, ingredient := range content.Ingredients {
		_, err = tx.Exec(insertIngredientsQuery, id, ingredient.IngredientID, ingredient.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *jsonMenuRepository) DeleteMenuItem(id int) error {
	tx, err := r.newDB.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	query := `
DELETE FROM menu_items
WHERE product_id =$1;
	`
	result, err := tx.Exec(query, id)
	if err != nil {
		return err
	}
	oneRes, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if oneRes == 0 {
		return errors.New("id incorrect")
	}
	return nil
}

func (r *jsonMenuRepository) GetMenuRepo() ([]models.MenuItem, error) {
	var menu []models.MenuItem
	query := `
	SELECT
		mi.product_id,
		mi.name,
		mi.description,
		mi.price,
		mi.category,
		mi.allergens,
		mii.ingredient_id,
		mii.quantity
	FROM menu_items mi
	LEFT JOIN menu_item_ingredients mii ON mi.product_id = mii.product_id;
	`
	rows, err := r.newDB.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	menuMap := make(map[int]*models.MenuItem)
	for rows.Next() {
		var ingredientId sql.NullInt64
		var quantity sql.NullFloat64
		var productId int
		var name, description, category string
		var allergens []string
		var price float64
		err = rows.Scan(
			&productId,
			&name,
			&description,
			&price,
			&category,
			pq.Array(&allergens),
			&ingredientId,
			&quantity,
		)
		if err != nil {
			return nil, err
		}
		if _, exists := menuMap[productId]; !exists {
			menuMap[productId] = &models.MenuItem{
				ID:          productId,
				Name:        name,
				Description: description,
				Price:       price,
				Category:    category,
				Allergens:   allergens,
				Ingredients: []models.MenuItemIngredient{},
			}

			menuMap[productId].Ingredients = append(menuMap[productId].Ingredients, models.MenuItemIngredient{
				IngredientID: int(ingredientId.Int64),
				Quantity:     float64(quantity.Float64),
			})

		}
	}
	for _, item := range menuMap {
		menu = append(menu, *item)
	}
	return menu, nil
}

func (r *jsonMenuRepository) GetMenuItemID(id int) (models.MenuItem, error) {
	var menuItem models.MenuItem
	query := `
	SELECT
		mi.product_id,
		mi.name,
		mi.description,
		mi.price,
		mi.category,
		mi.allergens,
		mii.ingredient_id,
		mii.quantity
	FROM menu_items mi
	LEFT JOIN menu_item_ingredients mii ON mi.product_id = mii.product_id
	WHERE mi.product_id = $1;
	`

	rows, err := r.newDB.Db.Query(query, id)
	if err != nil {
		return menuItem, err
	}
	found := false
	for rows.Next() {
		var (
			ingredientsID int
			quantity      float64
		)
		var (
			productId   int
			name        string
			description string
			price       float64
			category    string
			allergens   []string
		)
		err := rows.Scan(
			&productId,
			&name,
			&description,
			&price,
			&category,
			pq.Array(&allergens),
			&ingredientsID,
			&quantity,
		)
		if err != nil {
			return menuItem, err
		}
		if !found {
			menuItem = models.MenuItem{
				ID:          productId,
				Name:        name,
				Description: description,
				Price:       price,
				Category:    category,
				Allergens:   allergens,
				Ingredients: []models.MenuItemIngredient{},
			}
			found = true
		}
		if ingredientsID != 0 {
			menuItem.Ingredients = append(menuItem.Ingredients, models.MenuItemIngredient{
				IngredientID: ingredientsID,
				Quantity:     quantity,
			})
		}
	}
	if !found {
		return menuItem, fmt.Errorf("menu item with ID %d not found", id)
	}
	return menuItem, nil
}
