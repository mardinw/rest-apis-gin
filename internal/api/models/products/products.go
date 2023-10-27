package products

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"payuoge.com/internal/api/models"
)

type Product struct {
	ID           int64  `json:"id"`
	UserID       string `json:"user_id"`
	ProductCode  string `json:"product_code,omitempty"`
	ProductName  string `json:"product_name,omitempty"`
	Picture      string `json:"picture,omitempty"`
	Quantity     int16  `json:"quantity,omitempty"`
	Position     string `json:"position,omitempty"`
	SizeTypeId   int    `json:"size_type_id,omitempty"`
	SizeTypeName string `json:"size_type_name,omitempty"`
	CategoryId   int    `json:"category_id,omitempty"`
	CategoryName string `json:"category_name,omitempty"`
	BuyPrice     int32  `json:"buy_price,omitempty"`
	MRP          int32  `json:"min_retail_price,omitempty"`
	Total        int32  `json:"total_price"`
	Defective    int32  `json:"defective,omitempty"`
	Active       bool   `json:"active"`
	Created      int64  `json:"created,omitempty"`
	Updated      int64  `json:"updated,omitempty"`
}

func (product *Product) Insert(db *sql.DB, sizeTypeId, categoryId int, userId string) error {
	query := `
    INSERT INTO products(
    user_id,
	product_code,
	product_name,
	picture,
	quantity,
	position,
    size_type_id,
    category_id,
	buy_price, 
	mrp, 
	defective, 
	active, 
	created,
	updated)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    `

	times := time.Now().UnixMilli()
	args := []interface{}{
		userId,
		product.ProductCode,
		product.ProductName,
		product.Picture,
		product.Quantity,
		product.Position,
		sizeTypeId,
		categoryId,
		product.BuyPrice,
		product.MRP,
		product.Defective,
		product.Active,
		times,
		times,
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

func (product *Product) GetAll(db *sql.DB) ([]Product, error) {
	query := `
    SELECT 
    p.id,
	p.user_id,
    p.product_code,
    p.product_name,
    p.picture,
    p.quantity,
	p.position,
    s.name,
    c.name,
	p.mrp,
	p.buy_price,
	p.defective,
	p.active,
	p.created,
	p.updated
    FROM products p
	INNER JOIN category_products c ON p.category_id = c.id
	INNER JOIN size_type s ON p.size_type_id = s.id
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var result []Product
	for rows.Next() {
		var each = Product{}
		var err = rows.Scan(
			&each.ID,
			&each.UserID,
			&each.ProductCode,
			&each.ProductName,
			&each.Picture,
			&each.Quantity,
			&each.Position,
			&each.SizeTypeName,
			&each.CategoryName,
			&each.MRP,
			&each.BuyPrice,
			&each.Defective,
			&each.Active,
			&each.Created,
			&each.Updated,
		)
		if err != nil {
			log.Println(models.ErrRecordNotFound)
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil
}

func (product Product) GetProductsGroceries(userID string, db *sql.DB) ([]Product, error) {
	query := `
    SELECT 
    p.id,
	p.user_id,
    p.product_code,
    p.product_name,
    p.picture,
    p.quantity,
	p.position,
    s.name,
    c.name,
	p.mrp,
	p.buy_price,
	p.defective,
	p.active,
	p.created,
	p.updated
    FROM products p
	INNER JOIN category_products c ON p.category_id = c.id
	INNER JOIN size_type s ON p.size_type_id = s.id
    WHERE p.user_id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var result []Product
	for rows.Next() {
		var each = Product{}
		var err = rows.Scan(
			&each.ID,
			&each.UserID,
			&each.ProductCode,
			&each.ProductName,
			&each.Picture,
			&each.Quantity,
			&each.Position,
			&each.SizeTypeName,
			&each.CategoryName,
			&each.MRP,
			&each.BuyPrice,
			&each.Defective,
			&each.Active,
			&each.Created,
			&each.Updated,
		)
		if err != nil {
			log.Println(models.ErrRecordNotFound)
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil

}

func (product Product) Get(id int64, db *sql.DB) (*Product, error) {
	query := `
    SELECT 
    p.id,
	p.user_id,
    p.product_code,
    p.product_name,
    p.picture,
    p.quantity,
    p.position,
    s.name,
    c.name,
	p.mrp,
	p.buy_price,
	p.defective,
	p.active,
	p.created,
	p.updated
    FROM products p
	INNER JOIN category_products c ON p.category_id = c.id
	INNER JOIN size_type s ON p.size_type_id = s.id
	WHERE p.id = $1
	LIMIT 1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, query, id)

	if err := row.Scan(
		&product.ID,
		&product.UserID,
		&product.ProductCode,
		&product.ProductName,
		&product.Picture,
		&product.Quantity,
		&product.Position,
		&product.SizeTypeName,
		&product.CategoryName,
		&product.BuyPrice,
		&product.MRP,
		&product.Defective,
		&product.Active,
		&product.Created,
		&product.Updated,
	); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &product, nil
}

func (product *Product) Update(db *sql.DB, id int64, userId string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	query := `
    UPDATE products
    SET product_code = $1, 
    product_name = $2,
    picture = $3,
    quantity = $4,
    position = $5,
	size_type_id = $6,
	category_id = $7,
    mrp = $8,
    buy_price = $9,
    defective = $10,
	active = $11,
    updated = $12
    WHERE id = $13 AND user_id =$14
    `

	timeUpdate := time.Now().UnixMilli()

	args := []interface{}{
		product.ProductCode,
		product.ProductName,
		product.Picture,
		product.Quantity,
		product.Position,
		product.SizeTypeId,
		product.CategoryId,
		product.MRP,
		product.BuyPrice,
		product.Defective,
		product.Active,
		timeUpdate,
		id,
		userId,
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

func (product *Product) Delete(id int64, userId string, db *sql.DB) error {
	query := `
	DELETE FROM products
	WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{
		id, userId,
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
