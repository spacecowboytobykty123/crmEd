package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Teachers      TeacherModel
	Users         UserModel
	Cabinets      CabinetModel
	Subscriptions SubModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Teachers:      TeacherModel{DB: db},
		Users:         UserModel{DB: db},
		Cabinets:      CabinetModel{DB: db},
		Subscriptions: SubModel{DB: db},
	}
}
