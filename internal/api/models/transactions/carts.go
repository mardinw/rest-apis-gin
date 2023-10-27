package transactions

import (
	"database/sql"
	"log"
	"time"

	"golang.org/x/net/context"
	"payuoge.com/internal/api/models"
)

type Carts struct {
	ID           int64  `json:"id"`
	CustomerID   string `json:"customer_id"`
	ProductID    int64  `json:"product_id,omitempty"`
	ProductName  string `json:"product_name"`
	SizeTypeID   int8   `json:"size_type_id,omitempty"`
	SizeTypeName string `json:"size_type_name"`
	Quantity     int32  `json:"quantity"`
	Comments     string `json:"comments"`
	CreatedAt    int64  `json:"created_at"`
}

func (cart *Carts) Insert(userID string, db *sql.DB) error {
	query := `
    INSERT INTO carts(
    customer_id,
    product_id,
    quantity,
	size_type_id,
    comments,
    created_at
    ) VALUES($1, $2, $3, $4, $5,$6)
    `

	timeNow := time.Now().UnixMilli()
	args := []interface{}{
		userID,
		cart.ProductID,
		cart.Quantity,
		cart.SizeTypeID,
		cart.Comments,
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

func (carts *Carts) GetAll(userID string, db *sql.DB) ([]Carts, error) {
	query := `
    SELECT
    c.id,
    c.customer_id,
    p.product_name,
    c.quantity,
    s.name,
    c.comments,
    c.created_at
    FROM carts c
    INNER JOIN products p ON c.product_id = p.id
    INNER JOIN size_type s ON c.size_type_id = s.id
    WHERE c.customer_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []Carts
	for rows.Next() {
		var each = Carts{}
		var err = rows.Scan(
			&each.ID,
			&each.CustomerID,
			&each.ProductName,
			&each.Quantity,
			&each.SizeTypeName,
			&each.Comments,
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

func (cart *Carts) GetID(id int64, userID string, db *sql.DB) (*Carts, error) {
	query := `
    SELECT
    c.id,
    c.customer_id,
    p.product_name,
    c.quantity,
    s.name,
    c.comments,
    c.created_at
    FROM carts c
    INNER JOIN products p ON c.product_id = p.id
    INNER JOIN size_type s ON c.size_type_id = s.id
    WHERE c.id = $1 AND c.customer_id = $2
    LIMIT 1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	args := []interface{}{
		id,
		userID,
	}
	row := db.QueryRowContext(ctx, query, args...)

	if err := row.Scan(
		&cart.ID,
		&cart.CustomerID,
		&cart.ProductName,
		&cart.Quantity,
		&cart.SizeTypeName,
		&cart.Comments,
		&cart.CreatedAt,
	); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return cart, nil
}

func (cart *Carts) Update(id int64, userID string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	query := `
    UPDATE carts
	SET product_id = $1,
	quantity = $2,
    size_type_id = $3,
	comments = $4
	WHERE id = $5 AND customer_id $6
    `

	args := []interface{}{
		cart.ProductID,
		cart.Quantity,
		cart.SizeTypeID,
		cart.Comments,
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

func (cart *Carts) Delete(id int64, userID string, db *sql.DB) error {
	query := `
    DELETE FROM carts
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
