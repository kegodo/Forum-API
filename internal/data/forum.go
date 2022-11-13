// File: forum/internal/data/forum.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"forum.kevin.net/internal/validator"
)

// Todo struct supports the infromation for the todo todo
type Forum struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Publisher   string    `json:"publisher"`
	ReleaseDate int       `json:"releasedate"`
	Version     int32     `json:"version"`
}

func ValidateForum(v *validator.Validator, forum *Forum) {
	//using check() method to check our validation checks
	v.Check(forum.Title != "", "title", "must be provided")
	v.Check(len(forum.Title) <= 200, "title", "must not be more than 200 bytes long")

	v.Check(forum.Category != "", "category", "must be provided")
	v.Check(len(forum.Title) <= 250, "category", "must not be more than 250 bytes long")

	v.Check(forum.Description != "", "description", "must be provided")
	v.Check(len(forum.Description) <= 500, "description", "must not be more than 500 bytes long")

	v.Check(forum.Publisher != "", "Publisher", "must be provided")
	v.Check(len(forum.Publisher) <= 200, "Publisher", "must not be more than 200 bytes long")

}

type ForumModel struct {
	DB *sql.DB
}

// Insert() allows us to create a new todo
func (m ForumModel) Insert(forum *Forum) error {
	query := `
		INSERT INTO forums (title, category, description, publisher, releasedate)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, createdat, version
	`
	//collect the date field into a slice
	args := []interface{}{forum.Title, forum.Category, forum.Description, forum.Publisher, forum.ReleaseDate}
	//creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Clean up to prevent memory leaks
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&forum.ID, &forum.CreatedAt, &forum.Version)
}

// Get() allows us to retrieve a specific task
func (m ForumModel) Get(id int64) (*Forum, error) {
	//Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//Construct our query with the given id
	query := `
		SELECT id, createdat, title, category, description, publisher, releasedate, version
		FROM forums
		WHERE id = $1
	`

	//Declaring the Todo varaible to hold the returned data
	var forum Forum

	//Creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleaning up to prevent memory leaks
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&forum.ID,
		&forum.CreatedAt,
		&forum.Title,
		&forum.Category,
		&forum.Description,
		&forum.Publisher,
		&forum.ReleaseDate,
		&forum.Version,
	)

	if err != nil {
		//Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//Succes
	return &forum, nil
}

// Update() allows us to edit/alter a specific todo task
// Optimistic locking (version number)
func (m ForumModel) Update(forum *Forum) error {
	//create a query
	query := `
		UPDATE forums
		SET title = $1, category = $2, description = $3, publisher = $4, releasedate = $5, version = version + 1
		WHERE id = $6
		AND version = $7
		RETURNING version
	`
	args := []interface{}{
		forum.Title,
		forum.Category,
		forum.Description,
		forum.Publisher,
		forum.ReleaseDate,
		forum.ID,
		forum.Version,
	}

	//Creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleaning up to prevent memory leaks
	defer cancel()

	//Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&forum.Version)
	if err != nil {
		//Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	//Succes
	return nil

}

func (m ForumModel) Delete(id int64) error {
	//Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	//Create the delete query
	query := `
		DELETE FROM forums
		WHERE id = $1
	`

	//creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//clearing up to prevent memory leaks
	defer cancel()

	//Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	//Check how many rows were affected by the delete operation.
	//We call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	//Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m ForumModel) GetAll(title string, category string, description string, filters Filters) ([]*Forum, Metadata, error) {
	//constructing the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(),
	    id, createdat, title, category, description, publisher, releasedate, version
		FROM forums
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', category) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $3) OR $3 = '')
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortOrder())

	//creating the 3 second time out context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//Execute the query
	args := []interface{}{title, category, description, filters.limit(), filters.offSet()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	//Closing the result set
	defer rows.Close()
	totalRecords := 0

	//Initialize an empty slice to hold the task data
	forums := []*Forum{}

	//Iterate over the rows in the result set
	for rows.Next() {
		var forum Forum

		//Scanning the valus from the row into the todo struct
		err := rows.Scan(
			&totalRecords,
			&forum.ID,
			&forum.CreatedAt,
			&forum.Title,
			&forum.Category,
			&forum.Description,
			&forum.Publisher,
			&forum.ReleaseDate,
			&forum.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		//Add the todo to our slice
		forums = append(forums, &forum)
	}
	//checking for errors after looping through the result set
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	//returning the slice of todos
	return forums, metadata, nil
}
