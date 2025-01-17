package orderRepo

import (
	"errors"
	"fmt"
	"log"
)

func (r orderRepository) OrderClose(id int) error {
	tx, err := r.newDB.Db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic occurred: %v", p)

			tx.Rollback()
		} else if err != nil {
			log.Printf("Transaction rollback due to error: %v", err)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	stmt := `
UPDATE orders
SET status = 'close'
WHERE order_id = $1 AND status != 'close';
`
	res, err := tx.Exec(stmt, id)
	if err != nil {
		return err
	}
	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return errors.New("order not found or already closed")
	}
	return nil
}
