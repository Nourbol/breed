package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/breeds", app.createBreedHandler)
	router.HandlerFunc(http.MethodGet, "/v1/breeds/:id", app.showBreedHandler)
	router.HandlerFunc(http.MethodGet, "/v1/breeds", app.listBreedsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/breeds/:id", app.updateBreedHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/breeds/:id", app.deleteBreedHandler)

	return app.recoverPanic(router)
}
