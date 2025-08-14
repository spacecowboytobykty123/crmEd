package main

import (
	"authCRM/internal/data"
	"authCRM/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) createTeacherHandler(w http.ResponseWriter, r *http.Request) {
	var teacherinput struct {
		FullName  string             `json:"full_name"`
		BirthDate time.Time          `json:"birth_date"`
		Phone     string             `json:"phone"`
		Note      string             `json:"note"`
		Status    data.TeacherStatus `json:"status"`
		Gender    data.Gender        `json:"gender"`
	}

	err := app.readJSON(w, r, &teacherinput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	teacher := &data.Teacher{
		FullName:  teacherinput.FullName,
		BirthDate: teacherinput.BirthDate,
		Phone:     teacherinput.Phone,
		Note:      teacherinput.Note,
		Status:    data.StatusActive,
		Gender:    teacherinput.Gender,
	}

	v := validator.New()

	if data.ValidateTeacher(v, teacher); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Teachers.InsertTeacher(teacher)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/teachers/%d", teacher.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"teacher": teacher}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getTeacherHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	teacher, err := app.models.Teachers.GetTeacher(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"teacher": teacher}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	teacher, err := app.models.Teachers.GetTeacher(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var teacherinput struct {
		FullName  *string             `json:"full_name"`
		BirthDate *time.Time          `json:"birth_date"`
		Phone     *string             `json:"phone"`
		Note      *string             `json:"note"`
		Status    *data.TeacherStatus `json:"status"`
		Gender    *data.Gender        `json:"gender"`
	}

	err = app.readJSON(w, r, &teacherinput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if teacherinput.FullName != nil {
		teacher.FullName = *teacherinput.FullName
	}

	if teacherinput.BirthDate != nil {
		teacher.BirthDate = *teacherinput.BirthDate
	}

	if teacherinput.Phone != nil {
		teacher.Phone = *teacherinput.Phone
	}

	if teacherinput.Note != nil {
		teacher.Note = *teacherinput.Note
	}

	if teacherinput.Status != nil {
		teacher.Status = *teacherinput.Status
	}

	if teacherinput.Gender != nil {
		teacher.Gender = *teacherinput.Gender
	}

	v := validator.New()

	if data.ValidateTeacher(v, teacher); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Teachers.UpdateTeacher(teacher)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"teacher": teacher}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.models.Teachers.DeleteTeacher(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "успешно удалено"}, nil)
}

func (app *application) listTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var teacherInput struct {
		FullName      string
		TeacherStatus data.TeacherStatus
		Gender        data.Gender
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	teacherInput.FullName = app.readString(qs, "name", "")
	teacherInput.TeacherStatus = app.readTeacherStatus(qs, "status", "")
	teacherInput.Gender = app.readGender(qs, "gender", "")

	teacherInput.Filters.Page = app.readInt(qs, "page", 1, v)
	teacherInput.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	teacherInput.Filters.Sort = app.readString(qs, "sort", "id")
	teacherInput.Filters.SortSafelist = []string{"id", "name", "gender", "status", "-id", "-name", "-gender", "-status"}

	if data.ValidateFilters(v, teacherInput.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var gender *data.Gender
	if teacherInput.Gender != "" {
		gender = &teacherInput.Gender
	}

	var status *data.TeacherStatus
	if teacherInput.TeacherStatus != "" {
		status = &teacherInput.TeacherStatus
	}

	teachers, metadata, err := app.models.Teachers.GetAllTeachers(teacherInput.FullName, gender, status, teacherInput.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Teachers": teachers, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
