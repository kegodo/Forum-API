// Filename: cmd/api/context.go
package main

import (
	"context"
	"net/http"

	"forum.kevin.net/internal/data"
)

// Define a custom contextKey type
type contextKey string

// user key
const userContextKey = contextKey("user")

// Add user to context
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// Retrieve a user struct
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in the request context")
	}
	return user
}
