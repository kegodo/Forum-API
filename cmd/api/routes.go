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
	router.HandlerFunc(http.MethodGet, "/v1/forum", app.requirePermission("forum:read", app.listForumHandler))
	router.HandlerFunc(http.MethodPost, "/v1/forum", app.requirePermission("forum:write", app.createForumHandler))
	router.HandlerFunc(http.MethodGet, "/v1/forum/:id", app.requirePermission("forum:read", app.showForumHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/forum/:id", app.requirePermission("forum:write", app.updateForumHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/forum/:id", app.requirePermission("forum:write", app.deleteForumHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authentication(router))))
}
