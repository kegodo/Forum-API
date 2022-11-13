// File: forum/cmd/api/handlers.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"forum.kevin.net/internal/data"
	"forum.kevin.net/internal/validator"
)

func (app *application) createForumHandler(w http.ResponseWriter, r *http.Request) {
	//Our target decode destination
	var input struct {
		Title       string `json:"title"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Publisher   string `json:"publisher"`
		ReleaseDate int    `json:"releasedate"`
	}

	//Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//coping the valeus from the input struct to the new todo struct
	forum := &data.Forum{
		Title:       input.Title,
		Category:    input.Category,
		Description: input.Description,
		Publisher:   input.Publisher,
		ReleaseDate: input.ReleaseDate,
	}

	//Initialize a new Validator Instance
	v := validator.New()

	//check the map to determine if ther were any validation errors
	if data.ValidateForum(v, forum); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Creating a todo element
	err = app.models.Forums.Insert(forum)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	//Create a location header for the newly created resource
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/forums/%d", forum.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"forum": forum}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The showentry handler will display an individual todo element
func (app *application) showForumHandler(w http.ResponseWriter, r *http.Request) {
	//getting the request data from param function
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	//Fetching the specific todo element
	forum, err := app.models.Forums.Get(id)

	//Handling errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Writing the data from the returned get()
	err = app.writeJSON(w, http.StatusOK, envelope{"forum": forum}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Facilitates an update action to the todo element in the database
func (app *application) updateForumHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	//Fetch the original record from the database
	forum, err := app.models.Forums.Get(id)

	//Handling the errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Creating an input struct to hold data read in from the client
	//Updating the input struct to use pointers because pointers have a default value of nil
	var input struct {
		Title       *string `json:"title"`
		Category    *string `json:"category"`
		Description *string `json:"description"`
		Publisher   *string `json:"publisher"`
		ReleaseDate *int    `json:"releasedate"`
	}

	//Initilizing a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//checking for any updates
	if input.Title != nil {
		forum.Title = *input.Title
	}

	if input.Title != nil {
		forum.Category = *input.Category
	}

	if input.Description != nil {
		forum.Description = *input.Description
	}

	if input.Publisher != nil {
		forum.Publisher = *input.Publisher
	}

	if input.ReleaseDate != nil {
		forum.ReleaseDate = *input.ReleaseDate
	}

	//Initilize a new Validator Instance
	v := validator.New()

	//Checking the map to determin if there were any validation errors
	if data.ValidateForum(v, forum); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Passing the updated todo element to the update() method
	err = app.models.Forums.Update(forum)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Writing the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"forum": forum}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// To facilitate deletion of a todo element
func (app *application) deleteForumHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	err = app.models.Forums.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Returning 200 status ok to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "forum element sucessfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The listTodo handler allows the client to see a listing of todo elements based on a set of criteria
func (app *application) listForumHandler(w http.ResponseWriter, r *http.Request) {
	//creating an input struct to hold our query parameters
	var input struct {
		Title       string
		Category    string
		Description string
		data.Filters
	}

	//Initializing a validator
	v := validator.New()

	//getting the URL values map
	qs := r.URL.Query()

	//Using the helper method to extract the values
	input.Title = app.readString(qs, "title", "")
	input.Category = app.readString(qs, "category", "")
	input.Description = app.readString(qs, "decription", "")

	//Get the page information
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	//Get the sort information
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// Specific the allowed sort values
	input.Filters.SortList = []string{"id", "title", "cateogry", "description", "-id", "-title", "-category", "-description"}

	//checking for validation errors
	if data.ValidateFilter(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Geting a listing of all todo elements
	forums, metadata, err := app.models.Forums.GetAll(input.Title, input.Category, input.Description, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//sending JSON response
	err = app.writeJSON(w, http.StatusOK, envelope{"forums": forums, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
