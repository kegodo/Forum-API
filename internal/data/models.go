// File: forum/internal/data/models.go
package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// A wrapper for out data models
type Models struct {
	Forums ForumModel
	Users  UserModel
	Tokens TokenModel
}

// NewModels() allows us to create a new model
func NewModels(db *sql.DB) Models {
	return Models{
		Forums: ForumModel{DB: db},
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
	}
}
