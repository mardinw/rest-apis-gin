package transactions

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/net/context"
	"payuoge.com/internal/api/models"
)

type TransactionDetail struct {
	ID                int64  `json:"id"`
	TransactionID     int64  `json:"transaction_id"`
	ProductID         int64  `json:"product_id"`
	ProductName       string `json:"product_name"`
	PriceEach         int32  `json:"price_each"` // minimum retail price
	Quantity          int8   `json:"quantity"`
	SizeTypeID        int8   `json:"size_type_id"`
	SizeTypeName      string `json:"size_type"`
	TotalPriceProduct int32  `json:"total_price_product"`
	OrderLineNumber   int8   `json:"order_line_num"`
}

func (detail *TransactionDetail) Insert(transactionID int64, totalPrice int32, db *sql.DB) error {
	query := `
    INSERT INTO transaction_details(
    transaction_id,
    product_id,
    quantity,
    size_type_id,
    total,
    order_line_num
    ) VALUES($1,$2,$3,$4,$5,$6) 
    `
	args := []interface{}{
		transactionID,
		detail.ProductID,
		detail.Quantity,
		detail.SizeTypeID,
		totalPrice, // ini kita kalkulasikan setelah konfirm
		detail.OrderLineNumber,
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

func (detail *TransactionDetail) GetAll(db *sql.DB) ([]TransactionDetail, error) {
	query := `
    SELECT
    d.id,
    d.transaction_id,
    p.name,
    p.mrp,
    d.quantity,
    s.name,
    d.total,
    d.order_line_num
    from transaction_details d
    INNER JOIN products p on d.product_id = p.id
    INNER JOIN size_type s on d.size_type_id = s.id
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []TransactionDetail

	for rows.Next() {
		var each = TransactionDetail{}
		var err = rows.Scan(
			&each.ID,
			&each.TransactionID,
			&each.ProductName,
			&each.PriceEach,
			&each.Quantity,
			&each.SizeTypeName,
			&each.TotalPriceProduct,
			&each.OrderLineNumber,
		)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		result = append(result, each)

	}

	return result, nil
}

func (detail *TransactionDetail) Get(id int64, db *sql.DB) (*TransactionDetail, error) {
	query := `
    SELECT
    d.id,
    d.transaction_id,
    p.name,
    p.mrp,
    d.quantity,
    s.name,
    d.total,
    d.order_line_num
    from transaction_details d
    INNER JOIN products p on d.product_id = p.id
    INNER JOIN size_type s on d.size_type_id = s.id
    WHERE d.id = $1
    LIMIT 1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, query, id)

	if err := row.Scan(
		&detail.ID,
		&detail.TransactionID,
		&detail.ProductName,
		&detail.PriceEach,
		&detail.Quantity,
		&detail.SizeTypeName,
		&detail.TotalPriceProduct,
		&detail.OrderLineNumber,
	); err != nil {
		log.Println(err.Error())
		return nil, err

	}

	return detail, nil
}

func (detail *TransactionDetail) Update(id, transactionID int64, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	query := `
    UPDATE transaction_details
    SET
    product_id = $1,
    quantity = $2,
    size_type_id = $3,
    total = $4
    WHERE id = $5 AND transaction_id = $6
    `

	args := []interface{}{
		detail.ProductID,
		detail.Quantity,
		detail.SizeTypeID,
		detail.TotalPriceProduct,
		id,
		transactionID,
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

func (detail *TransactionDetail) Delete(id int64, db *sql.DB) error {
	query := `
    DELETE FROM transaction_details
    WHERE id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{
		id,
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
