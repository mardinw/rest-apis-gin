package transactions

import (
	"context"
	"database/sql"
	"log"
	"time"

	"payuoge.com/internal/api/models"
	"payuoge.com/internal/api/models/products"
)

type Checkouts struct {
	ID          int64   `json:"id"`
	CustomerID  string  `json:"customer_id"`
	TotalAmount float64 `json:"total_amount"`
	CreatedAt   int64   `json:"created_at"`
}

func (checkout *Checkouts) CalculateTotalAmount(userID string, db *sql.DB) float64 {
	var carts Carts
	var product products.Product

	query := `
    SELECT
    c.product_id,
    c.quantity,
    p.mrp
    FROM carts c
    INNER JOIN products p ON c.product_id = p.id
    WHERE c.customer_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Printf("Error querying cart items: %v", err)
		return 0.0
	}

	defer rows.Close()
	totalAmount := 0.0
	for rows.Next() {
		if err := rows.Scan(
			&carts.ProductID,
			&carts.Quantity,
			&product.MRP,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		totalAmount += float64(carts.Quantity) * float64(product.MRP)

	}

	return totalAmount
}

func (checkout *Checkouts) Insert(userID string, totalAmount float64, db *sql.DB) error {
	query := `
    INSERT INTO checkouts(
    customer_id,
    total_amount,
    created_at
    ) VALUES($1, $2, $3)
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

func (checkout *Checkouts) GetAll(userID string, db *sql.DB) ([]Checkouts, error) {
	query := `
    SELECT * FROM checkouts WHERE customer_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []Checkouts
	for rows.Next() {
		var each = Checkouts{}
		var err = rows.Scan(
			&each.ID,
			&each.CustomerID,
			&each.TotalAmount,
			&each.CreatedAt,
		)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil
}

func (checkout *Checkouts) GetID(id int64, userID string, db *sql.DB) (*Checkouts, error) {
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
		&checkout.ID,
		&checkout.CustomerID,
		&checkout.TotalAmount,
		&checkout.CreatedAt,
	); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return checkout, nil
}

func (checkout *Checkouts) Update(id int64, userID string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	query := `
    UPDATE checkouts
    SET total_amount = $1
    WHERE id =$2 AND customer_id = $3
    `

	args := []interface{}{
		checkout.TotalAmount,
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

func (checkout *Checkouts) Delete(id int64, userID string, db *sql.DB) error {
	query := `
    DELETE FROM checkouts
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
