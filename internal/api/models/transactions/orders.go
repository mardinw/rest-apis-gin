package transactions

import (
	"context"
	"database/sql"
	"log"
	"time"

	"payuoge.com/internal/api/models"
)

type Orders struct {
	ID          int64  `json:"id"`
	CustomerID  string `json:"customer_id"`
	TotalAmount int32  `json:"total_amount"`
	OrderDate   int64  `json:"order_date"`
}

func (order *Orders) Insert(userID string, totalAmount int32, db *sql.DB) error {
	query := `
    INSERT INTO orders (
    customer_id,
    total_amount,
    order_date
    ) VALUES ($1, $2, $3)
    `
	timeNow := time.Now().UnixMilli()

	args := []interface{}{
		userID,
		totalAmount,
		timeNow,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (order *Orders) GetAll(userID string, db *sql.DB) ([]Orders, error) {
	query := `
    SELECT * FROM orders WHERE customer_id = $1 
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var result []Orders
	for rows.Next() {
		var each = Orders{}
		var err = rows.Scan(
			&each.ID,
			&each.CustomerID,
			&each.TotalAmount,
			&each.OrderDate,
		)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil
}

func (order *Orders) GetID(id int64, userID string, db *sql.DB) (*Orders, error) {
	query := `
    SELECT * FROM checkouts WHERE id = $1 AND customer_id = $2 LIMIT 1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{
		id, userID,
	}

	row := db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.TotalAmount,
		&order.OrderDate,
	); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return order, nil
}

func (order *Orders) Update(id int64, userID string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	query := `
    UPDATE orders
    SET total_amount = $1
    WHERE id = $2 AND customer_id = $3
    `

	args := []interface{}{
		order.TotalAmount,
		id, userID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (order *Orders) Delete(id int64, userID string, db *sql.DB) error {
	query := `
    DELETE FROM orders
    WHERE id = $1 AND customer_id = $2
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	args := []interface{}{
		id, userID,
	}

	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if rowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil
}
