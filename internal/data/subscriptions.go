package data

import (
	"authCRM/internal/validator"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Price          int32     `json:"price"`
	Type           SubStatus `json:"type"`
	DurationMonths *int16    `json:"duration_months,omitempty"`
	SessionsCount  *int16    `json:"sessions_count,omitempty"`
	ValidityMonths *int16    `json:"validity_months,omitempty"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func getValue(v *int16) int16 {
	if v == nil {
		return 0
	}
	return *v
}

func ValidateSubscription(v *validator.Validator, sub *Subscription) {
	v.Check(sub.Name != "", "name", "должны добавить имя!")
	v.Check(len(sub.Name) <= 200, "name", "имя не больше 200 байтов!")
	v.Check(sub.Type != "", "type", "должны выбрать тип!")
	v.Check(sub.Price > 0, "price", "сумма должна быть больше нуля!")

	if sub.Type == Monthly {
		v.Check(getValue(sub.ValidityMonths) == 0, "validMonths", "такой параметр не для периодной подписки")
		v.Check(getValue(sub.SessionsCount) == 0, "sessionCount", "такой параметр не для периодной подписки")
	}
	if sub.Type == Visits {
		v.Check(getValue(sub.DurationMonths) == 0, "durationMonths", "такой параметр не для количественной подписки")
	}
}

type SubModel struct {
	DB *sql.DB
}

func (s SubModel) InsertSubscription(sub *Subscription) error {

	query := `INSERT INTO subscriptions (name, price, type, duration_months, sessions_count, validity_months)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at
`

	args := []any{sub.Name, sub.Price, sub.Type, sub.DurationMonths, sub.SessionsCount, sub.ValidityMonths}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&sub.ID, &sub.CreatedAt)
}

func (s SubModel) GetSubscription(id uuid.UUID) (*Subscription, error) {
	query := `SELECT id, name, price, type, duration_months, sessions_count, validity_months, updated_at
	FROM subscriptions
	WHERE id = $1
	`

	var sub Subscription

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&sub.ID,
		&sub.Name,
		&sub.Price,
		&sub.Type,
		&sub.DurationMonths,
		&sub.SessionsCount,
		&sub.ValidityMonths,
		&sub.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &sub, err
}

func (s SubModel) UpdateSubscription(sub *Subscription) error {
	query := `UPDATE subscriptions
	SET name = $1, price = $2, type = $3, duration_months = $4, sessions_count = $5, validity_months = $6, updated_at = NOW()
	WHERE id = $7 and updated_at = $8
	RETURNING updated_at
`

	updatedAt := sub.UpdatedAt.UTC().Truncate(time.Microsecond)

	args := []any{sub.Name, sub.Price, sub.Type, sub.DurationMonths, sub.SessionsCount, sub.ValidityMonths, sub.ID, updatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.DB.QueryRowContext(ctx, query, args...).Scan(&sub.UpdatedAt)
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

func (s SubModel) DeleteSubscription(id uuid.UUID) error {
	query := `DELETE FROM subscriptions
	WHERE id = $1

`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, query, id)
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

func (s SubModel) GetAllSubscriptions() ([]*Subscription, error) {
	query := `SELECT COUNT(*) OVER(), id, name, price, type FROM subscriptions`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	totalRecords := 0
	subs := []*Subscription{}

	for rows.Next() {
		var sub Subscription

		err := rows.Scan(
			&totalRecords,
			&sub.ID,
			&sub.Name,
			&sub.Price,
			&sub.Type,
		)
		if err != nil {
			return nil, err
		}

		subs = append(subs, &sub)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}
