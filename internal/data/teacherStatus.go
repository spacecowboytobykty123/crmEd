package data

type TeacherStatus string

const (
	StatusActive   TeacherStatus = "активный"
	StatusVacation TeacherStatus = "отпуск"
	StatusArchived TeacherStatus = "архивный"
)
