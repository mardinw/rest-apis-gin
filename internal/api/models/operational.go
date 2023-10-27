package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lib/pq"
)

type Operationals struct {
	ID             uint64         `json:"id"`
	GroceriesID    string         `json:"groceries_id"`
	DayOperational pq.StringArray `json:"day_operational,omitempty"`
	Open           int64          `json:"open,omitempty"`
	Close          int64          `json:"close,omitempty"`
	Active         bool           `json:"active"`
}

func (operate *Operationals) Insert(userId string, db *sql.DB) error {
	query := `
    INSERT INTO operationals(
    groceries_id,
    day_operational,
    open,
    close,
    active
    ) VALUES (
    $1, $2, $3, $4, $5
    )
    `

	dayOperational := pq.Array(operate.DayOperational)

	args := []interface{}{
		userId,
		dayOperational,
		operate.Open,
		operate.Close,
		operate.Active,
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

func (operate *Operationals) GetAll(db *sql.DB) ([]Operationals, error) {
	query := `
    SELECT * FROM operationals
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var result []Operationals
	for rows.Next() {
		var each = Operationals{}
		var err = rows.Scan(
			&each.ID,
			&each.GroceriesID,
			&each.DayOperational,
			&each.Open,
			&each.Close,
			&each.Active,
		)

		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		result = append(result, each)
	}

	return result, nil
}

func (operate *Operationals) Get(id uint64, db *sql.DB) (*Operationals, error) {
	query := `
    SELECT * FROM operationals
    WHERE id = $1
    LIMIT 1
    `

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, query, id)

	if err := row.Scan(
		&operate.ID,
		&operate.GroceriesID,
		&operate.DayOperational,
		&operate.Open,
		&operate.Close,
		&operate.Active,
	); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return operate, nil
}

func (operate *Operationals) Update(id uint64, userId string, db *sql.DB) error {
	query := `
    UPDATE operationals
    SET open = $1,
    close = $2,
    active = $3
    WHERE id = $4 AND groceries_id = $5
    RETURNING id
    `

	args := []interface{}{
		operate.Open,
		operate.Close,
		operate.Active,
		id,
		userId,
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

func (operate *Operationals) Delete(id uint64, userId string, db *sql.DB) error {
	query := `
    DELETE FROM operationals
    WHERE id = $1 AND groceries_id = $2
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
		return ErrRecordNotFound
	}

	return nil
}
