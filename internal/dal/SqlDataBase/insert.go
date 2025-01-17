package SqlDataBase

import (
	"log/slog"
	"os"
)

func (d *DB) InsertInto() error {
	tx, err := d.Db.Begin()
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
	tables := []string{
		"inventory",
		"menu_items",
		"orders",
		"order_items",
		"price_history",
		"order_status_history",
		"inventory_transactions",
		"menu_item_ingredients",
	}
	for _, table := range tables {
		query := "SELECT * FROM " + table
		res, err := tx.Exec(query)
		if err != nil {
			return err
		}
		rowsAff, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAff > 0 {
			slog.Info("tables already contain data, skipping insert")
			return nil
		}
	}

	content, err := os.ReadFile("../insert.sql")
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(content))
	if err != nil {
		return err
	}

	return nil
}
