package main

import (
	"authCRM/internal/data"
	"authCRM/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createCabinetHandler(w http.ResponseWriter, r *http.Request) {
	var cabinetInput struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	err := app.readJSON(w, r, &cabinetInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cabinet := &data.Cabinet{
		Name:    cabinetInput.Name,
		Address: cabinetInput.Address,
	}

	v := validator.New()

	if data.ValidateCabinet(v, cabinet); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Cabinets.InsertCabinet(cabinet)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/teachers/%d", cabinet.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"teacher": cabinet}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCabinetHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	cabinet, err := app.models.Cabinets.GetCabinet(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"teacher": cabinet}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateCabinetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	cabinet, err := app.models.Cabinets.GetCabinet(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var cabinetinput struct {
		Name    *string `json:"name"`
		Address *string `json:"address"`
	}

	err = app.readJSON(w, r, &cabinetinput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if cabinetinput.Name != nil {
		cabinet.Name = *cabinetinput.Name
	}

	if cabinetinput.Address != nil {
		cabinet.Address = *cabinetinput.Address
	}

	v := validator.New()

	if data.ValidateCabinet(v, cabinet); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Cabinets.UpdateCabinet(cabinet)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cabinet": cabinet}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCabinetHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.models.Cabinets.DeleteCabinet(id)
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

func (app *application) listCabinetsHandler(w http.ResponseWriter, r *http.Request) {
	var cabinetInput struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	cabinetInput.Filters.Page = app.readInt(qs, "page", 1, v)
	cabinetInput.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	cabinetInput.Filters.Sort = app.readString(qs, "sort", "id")

	cabinetInput.Filters.SortSafelist = []string{"id", "-id"}
	if data.ValidateFilters(v, cabinetInput.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	cabinets, metadata, err := app.models.Cabinets.GetAllCabinets(cabinetInput.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Cabinets": cabinets, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
