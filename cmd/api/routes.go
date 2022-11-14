// File: forum/cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	//security routes
	router.NotFound = http.HandlerFunc(app.notFoundReponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/forum", app.listForumHandler)
	router.HandlerFunc(http.MethodPost, "/v1/forum", app.createForumHandler)
	router.HandlerFunc(http.MethodGet, "/v1/forum/:id", app.showForumHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/forum/:id", app.updateForumHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/forum/:id", app.deleteForumHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
