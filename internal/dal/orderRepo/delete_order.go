package orderRepo

func (r orderRepository) DeleteOrder(id int) error {
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
	stmt := `
	SELECT status FROM orders WHERE order_id= $1;
	`
	status := ""
	row := tx.QueryRow(stmt, id)
	err = row.Scan(&status)
	if err != nil {
		return err
	}
	err = row.Err()
	if err != nil {
		return err
	}
	if status == "open" {
		err = r.RestockInventory(tx, id)
		if err != nil {
			return err
		}
	}

	err = r.DeleteOldOrder(tx, id)
	if err != nil {
		return err
	}
	return nil
}
