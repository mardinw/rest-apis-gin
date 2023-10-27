package products

import (
	"context"
	"database/sql"
	"log"
	"time"

	"payuoge.com/internal/api/models"
)

type CategoryProducts struct {
	ID          int       `json:"id"`
	UserID      string    `json:"user_id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Products    []Product `json:"product,omitempty"`
}

func (category *CategoryProducts) Insert(db *sql.DB, userId string) error {
	query := `
    INSERT INTO category_products(
    name,
    user_id,
    Description
    )
    VALUES($1, $2, $3)
    `
	args := []interface{}{
		category.Name,
		userId,
		category.Description,
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

func (category *CategoryProducts) GetAll(db *sql.DB) ([]CategoryProducts, error) {
	query := `
    SELECT id, name, description FROM category_products
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []CategoryProducts
	for rows.Next() {
		var each = CategoryProducts{}
		var err = rows.Scan(
			&each.ID,
			&each.Name,
			&each.Description,
		)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		result = append(result, each)
	}
	return result, nil
}

func (category *CategoryProducts) Get(id int, db *sql.DB) (*CategoryProducts, error) {
	query := `
    SELECT id, name, description FROM category_products WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	args := []interface{}{
		id,
	}
	if err := db.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
	); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return category, nil
}

func (category *CategoryProducts) Update(id int, userId string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	query := `
    UPDATE category_products
    SET name = $1,
    description = $2
    WHERE id = $3 AND user_id = $4
    RETURNING id
    `

	args := []interface{}{
		category.Name,
		category.Description,
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
		log.Fatal(err)
		return err
	}

	return nil
}

func (category *CategoryProducts) Delete(id int, userId string, db *sql.DB) error {
	query := `
    DELETE FROM category_products
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
