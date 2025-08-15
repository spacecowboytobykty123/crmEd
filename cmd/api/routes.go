package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/teacher", app.createTeacherHandler)
	router.HandlerFunc(http.MethodGet, "/v1/teacher/:id", app.getTeacherHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/teacher/:id", app.updateTeacherHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/teacher/:id", app.deleteTeacherHandler)
	router.HandlerFunc(http.MethodGet, "/v1/teachers", app.listTeachersHandler)

	router.HandlerFunc(http.MethodPost, "/v1/user", app.registerUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/cabinet/:id", app.getCabinetHandler)
	router.HandlerFunc(http.MethodPost, "/v1/cabinet", app.createCabinetHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/cabinet/:id", app.updateCabinetHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/cabinet/:id", app.deleteCabinetHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cabinets", app.listCabinetsHandler)

	router.HandlerFunc(http.MethodGet, "/v1/subscription/:id", app.getSubHandler)
	router.HandlerFunc(http.MethodPost, "/v1/subscription/", app.createSubHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/subscription/:id", app.updateSubHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/subscription/:id", app.deleteSubscriptionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/subscriptions", app.listSubscriptionsHandler)

	return app.recoverPanic(app.rateLimit(router))
}
