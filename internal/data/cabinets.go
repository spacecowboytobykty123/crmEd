package data

import (
	"authCRM/internal/validator"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

type Cabinet struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Version int       `json:"-"`
}

func ValidateCabinet(v *validator.Validator, cabinet *Cabinet) {
	v.Check(cabinet.Name != "", "name", "Имя кабинета не должно быть пустым!")
}

type CabinetModel struct {
	DB *sql.DB
}

func (c CabinetModel) InsertCabinet(cabinet *Cabinet) error {
	query := `INSERT INTO cabinets (name, address)
	VALUES ($1, $2)
	RETURNING id, version
`

	args := []any{cabinet.Name, cabinet.Address}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&cabinet.ID, &cabinet.Version)
}

func (c CabinetModel) GetCabinet(id uuid.UUID) (*Cabinet, error) {
	query := `SELECT id, name, address, version
	FROM cabinets
	WHERE id = $1
`

	var cabinet Cabinet

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&cabinet.ID,
		&cabinet.Name,
		&cabinet.Address,
		&cabinet.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &cabinet, err
}

func (c CabinetModel) UpdateCabinet(cabinet *Cabinet) error {
	query := `UPDATE cabinets
	SET name = $1, address = $2, version = version + 1
	WHERE id = $3 and version = $4
	RETURNING version
`

	args := []any{cabinet.Name, cabinet.Address, cabinet.ID, cabinet.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&cabinet.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err

		}
	}
	return nil
}

func (c CabinetModel) DeleteCabinet(id uuid.UUID) error {
	query := `DELETE FROM cabinets
	WHERE id = $1

`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil

}

func (c CabinetModel) GetAllCabinets(filters Filters) ([]*Cabinet, Metadata, error) {
	query := `SELECT COUNT(*) OVER(), id, name, address FROM cabinets`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	cabinets := []*Cabinet{}

	for rows.Next() {
		var cabinet Cabinet

		err := rows.Scan(
			&totalRecords,
			&cabinet.ID,
			&cabinet.Name,
			&cabinet.Address,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		cabinets = append(cabinets, &cabinet)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return cabinets, metadata, nil
}
