package data

import (
	"authCRM/internal/validator"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

type Teacher struct {
	ID           uuid.UUID     `json:"id"`
	FullName     string        `json:"full_name"`
	BirthDate    time.Time     `json:"birth_date"`
	Phone        string        `json:"phone"`
	Note         string        `json:"note"`
	Gender       Gender        `json:"gender"`
	Status       TeacherStatus `json:"status"`
	CreatedAt    time.Time     `json:"-"`
	UpdatedAt    time.Time     `json:"-"`
	SalaryRateID int32         `json:"salary_rate_id"`
}

func ValidateTeacher(v *validator.Validator, teacher *Teacher) {
	v.Check(teacher.FullName != "", "name", "должны добавить имя!")
	v.Check(len(teacher.FullName) <= 200, "name", "имя не больше 200 байтов!")
	v.Check(teacher.Phone != "", "phone", "должны добавить телефон!")
}

type TeacherModel struct {
	DB *sql.DB
}

func (t TeacherModel) InsertTeacher(teacher *Teacher) error {
	query := `INSERT INTO teachers (full_name, birth_date, phone, note, status, gender)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at
`

	args := []any{teacher.FullName, teacher.BirthDate, teacher.Phone, teacher.Note, teacher.Status, teacher.Gender}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&teacher.ID, &teacher.CreatedAt)
}

func (t TeacherModel) GetTeacher(id uuid.UUID) (*Teacher, error) {
	query := `SELECT id, full_name, birth_date, phone, note, status, updated_at, gender
	FROM teachers
	WHERE id = $1
	`

	var teacher Teacher

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		&teacher.ID,
		&teacher.FullName,
		&teacher.BirthDate,
		&teacher.Phone,
		&teacher.Note,
		&teacher.Status,
		&teacher.UpdatedAt,
		&teacher.Gender,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &teacher, err
}

func (t TeacherModel) UpdateTeacher(teacher *Teacher) error {
	query := `UPDATE teachers
	SET full_name = $1, birth_date = $2, phone = $3, note = $4, status = $5, gender = $6,updated_at = NOW()
	WHERE id = $7 and updated_at = $8
	RETURNING updated_at
`

	args := []any{teacher.FullName, teacher.BirthDate, teacher.Phone, teacher.Note, teacher.Status, teacher.Gender, teacher.ID, teacher.UpdatedAt}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&teacher.UpdatedAt)
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

func (t TeacherModel) DeleteTeacher(id uuid.UUID) error {
	query := `DELETE FROM teachers
	WHERE id = $1

`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
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

func (t TeacherModel) GetAllTeachers(name string, gender *Gender, status *TeacherStatus, filters Filters) ([]*Teacher, Metadata, error) {
	query := `SELECT COUNT(*) OVER(), id, full_name, birth_date, phone, status, gender
FROM teachers
WHERE (to_tsvector('simple', full_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
  AND ($2::gender IS NULL OR gender = $2::gender)
  AND ($3::teacher_status IS NULL OR status = $3::teacher_status)
ORDER BY id
LIMIT $4 OFFSET $5`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, name, gender, status, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	teachers := []*Teacher{}

	for rows.Next() {
		var teacher Teacher

		err := rows.Scan(
			&totalRecords,
			&teacher.ID,
			&teacher.FullName,
			&teacher.BirthDate,
			&teacher.Phone,
			&teacher.Status,
			&teacher.Gender,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		teachers = append(teachers, &teacher)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return teachers, metadata, nil
}
