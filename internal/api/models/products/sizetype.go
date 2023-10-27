package products

import (
	"context"
	"database/sql"
	"log"
	"time"

	"payuoge.com/internal/api/models"
)

type SizeType struct {
	ID       int       `json:"id"`
	UserID   string    `json:"user_id,omitempty"`
	Name     string    `json:"name"`
	Products []Product `json:"product,omitempty"`
}

func (size *SizeType) Insert(db *sql.DB, userId string) error {
	query := `
    INSERT INTO size_type(
    user_id,
    name
    )
    VALUES ($1, $2)
    `
	args := []interface{}{
		userId,
		size.Name,
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

func (size *SizeType) GetAll(db *sql.DB) ([]SizeType, error) {
	query := `
    SELECT id, name FROM size_type
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var result []SizeType
	for rows.Next() {
		var each = SizeType{}
		var err = rows.Scan(
			&each.ID,
			&each.Name,
		)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil
}

func (size *SizeType) Get(id int, db *sql.DB) (*SizeType, error) {
	query := `
    SELECT id, name FROM size_type
    WHERE id = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{
		id,
	}
	if err := db.QueryRowContext(ctx, query, args...).Scan(
		&size.ID,
		&size.Name,
	); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return size, nil
}

func (size *SizeType) Update(id int, userId string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	query := `
    UPDATE size_type
    SET name = $1
    WHERE id = $2 AND user_id = $3
    RETURNING id
    `

	args := []interface{}{
		size.Name,
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
	}

	return nil
}

func (size *SizeType) Delete(id int, userId string, db *sql.DB) error {
	query := `
	DELETE FROM size_type
	WHERE id = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []interface{}{
		id,
		userId,
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
