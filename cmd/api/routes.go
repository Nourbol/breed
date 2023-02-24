package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.requirePermission("breeds: read", app.healthcheckHandler))
	router.HandlerFunc(http.MethodPost, "/v1/breeds", app.requirePermission("breeds: write", app.createBreedHandler))
	router.HandlerFunc(http.MethodGet, "/v1/breeds/:id", app.requirePermission("breeds: read", app.showBreedHandler))
	router.HandlerFunc(http.MethodGet, "/v1/breeds", app.requirePermission("breeds: read", app.listBreedsHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/breeds/:id", app.requirePermission("breeds: write", app.updateBreedHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/breeds/:id", app.requirePermission("breeds: write", app.deleteBreedHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimiter(app.authenticate(router)))
}
