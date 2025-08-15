package main

import (
	"authCRM/internal/data"
	"authCRM/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createSubHandler(w http.ResponseWriter, r *http.Request) {
	var subInput struct {
		Name           string         `json:"name"`
		Price          int32          `json:"price"`
		Type           data.SubStatus `json:"type"`
		DurationMonths *int16         `json:"duration_months,omitempty"`
		SessionsCount  *int16         `json:"sessions_count,omitempty"`
		ValidityMonths *int16         `json:"validity_months,omitempty"`
	}

	err := app.readJSON(w, r, &subInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	sub := &data.Subscription{
		Name:           subInput.Name,
		Price:          subInput.Price,
		Type:           subInput.Type,
		DurationMonths: subInput.DurationMonths,
		SessionsCount:  subInput.SessionsCount,
		ValidityMonths: subInput.ValidityMonths,
	}

	v := validator.New()

	if data.ValidateSubscription(v, sub); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Subscriptions.InsertSubscription(sub)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/teachers/%d", sub.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"subscription": sub}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSubHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	sub, err := app.models.Subscriptions.GetSubscription(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"subscription": sub}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateSubHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	sub, err := app.models.Subscriptions.GetSubscription(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var subinput struct {
		Name           *string         `json:"name"`
		Price          *int32          `json:"price"`
		Type           *data.SubStatus `json:"type"`
		DurationMonths *int16          `json:"duration_months,omitempty"`
		SessionsCount  *int16          `json:"sessions_count,omitempty"`
		ValidityMonths *int16          `json:"validity_months,omitempty"`
	}

	err = app.readJSON(w, r, &subinput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if subinput.Name != nil {
		sub.Name = *subinput.Name
	}

	if subinput.Price != nil {
		sub.Price = *subinput.Price
	}

	if subinput.Type != nil {
		sub.Type = *subinput.Type
	}

	if subinput.DurationMonths != nil {
		sub.DurationMonths = subinput.DurationMonths
	}

	if subinput.SessionsCount != nil {
		sub.SessionsCount = subinput.SessionsCount
	}

	if subinput.ValidityMonths != nil {
		sub.ValidityMonths = subinput.ValidityMonths
	}

	v := validator.New()

	if data.ValidateSubscription(v, sub); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Subscriptions.UpdateSubscription(sub)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"subscription": sub}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	err = app.models.Subscriptions.DeleteSubscription(id)
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

func (app *application) listSubscriptionsHandler(w http.ResponseWriter, r *http.Request) {

	subs, err := app.models.Subscriptions.GetAllSubscriptions()
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"subscriptions": subs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
