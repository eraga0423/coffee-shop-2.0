package dal

import (
	"frapuccino/internal/dal/SqlDataBase"
	"frapuccino/models"
)

// InventoryRepository defines the methods for reading and writing inventory data.
type InventoryRepository interface {
	ReadJSONInv() ([]models.InventoryItem, error)   // Reads the inventory data from a JSON file.
	WriteJSONInv(body []models.InventoryItem) error // Writes the updated inventory data to a JSON file.
	AddItems(item models.InventoryItem) error
	UpdateItem(int, models.InventoryItem) error
	DeleteItem(int) error
	CheckIfExists(ingredientID int) (bool, error)
	CheckIfNameExists(name string) (bool, error)
}

// jsonInvRepository implements the InventoryRepository interface using JSON file storage.
type jsonInvRepository struct {
	newDB *SqlDataBase.DB
}

// NewJSONInvRepository creates and returns a new instance of jsonInvRepository.
func NewJSONInvRepository(db *SqlDataBase.DB) InventoryRepository {
	return &jsonInvRepository{newDB: db}
}

func (j *jsonInvRepository) ReadJSONInv() ([]models.InventoryItem, error) {
	rows, err := j.newDB.Db.Query(`SELECT ingredient_id, name, quantity, unit FROM inventory`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		err := rows.Scan(&item.IngredientID, &item.Name, &item.Quantity, &item.Unit)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *jsonInvRepository) WriteJSONInv(newInventory []models.InventoryItem) error {
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

	query := `INSERT INTO inventory (name, quantity, unit) VALUES ($1, $2, $3)`
	for _, item := range newInventory {
		_, err := r.newDB.Db.Exec(query, item.Name, item.Quantity, item.Unit)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *jsonInvRepository) AddItems(item models.InventoryItem) error {
	tx, err := j.newDB.Db.Begin()
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

	query := `INSERT INTO inventory (name, quantity, unit) VALUES ($1, $2, $3)`
	_, err1 := tx.Exec(query, item.Name, item.Quantity, item.Unit)
	if err1 != nil {
		return err
	}
	return nil
}

func (j *jsonInvRepository) UpdateItem(id int, item models.InventoryItem) error {
	tx, err := j.newDB.Db.Begin()
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

	query := `UPDATE inventory SET name = $1, quantity = $2, unit = $3 WHERE ingredient_id = $4`
	_, err1 := j.newDB.Db.Exec(query, item.Name, item.Quantity, item.Unit, id)
	if err1 != nil {
		return err1
	}
	return nil
}

func (j *jsonInvRepository) DeleteItem(id int) error {
	tx, err := j.newDB.Db.Begin()
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
	query := `DELETE FROM inventory WHERE ingredient_id = $1`
	_, err1 := j.newDB.Db.Exec(query, id)
	if err1 != nil {
		return err1
	}
	return nil
}

func (j *jsonInvRepository) CheckIfExists(ingredientID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM inventory WHERE ingredient_id = $1)`
	err := j.newDB.Db.QueryRow(query, ingredientID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (j *jsonInvRepository) CheckIfNameExists(name string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM inventory WHERE name = $1)`
	err := j.newDB.Db.QueryRow(query, name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
