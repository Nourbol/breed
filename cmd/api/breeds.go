package main

import (
	"errors"
	"fmt"
	"github.com/Nourbol/breed/internal/data"
	"github.com/Nourbol/breed/internal/validator"
	"net/http"
)

func (app *application) createBreedHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		AvgCost     data.Cost `json:"avg_cost"`
		Countries   []string  `json:"countries"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	breed := &data.Breed{
		Name:        input.Name,
		Description: input.Description,
		AvgCost:     input.AvgCost,
		Countries:   input.Countries,
	}

	v := validator.New()
	if data.ValidateBreed(v, breed); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err = app.models.Breeds.Insert(breed); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/breeds/%d", breed.ID))

	if err = app.writeJSON(w, http.StatusCreated, envelope{"breed": breed}, headers); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showBreedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	breed, err := app.models.Breeds.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"breed": breed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateBreedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	breed, err := app.models.Breeds.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name        *string    `json:"name"`
		Description *string    `json:"description"`
		AvgCost     *data.Cost `json:"avg_cost"`
		Countries   []string   `json:"countries"`
	}

	if err = app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		breed.Name = *input.Name
	}
	if input.Description != nil {
		breed.Description = *input.Description
	}
	if input.AvgCost != nil {
		breed.AvgCost = *input.AvgCost
	}
	if input.Countries != nil {
		breed.Countries = input.Countries
	}

	v := validator.New()
	if data.ValidateBreed(v, breed); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err = app.models.Breeds.Update(breed); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err = app.writeJSON(w, http.StatusOK, envelope{"breed": breed}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBreedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if err = app.models.Breeds.Delete(id); err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err = app.writeJSON(w, http.StatusOK, envelope{"message": "breed successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listBreedsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string
		Countries []string
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Name = app.readStringQueryParam(qs, "name", "")
	input.Countries = app.readCSVQueryParam(qs, "countries", []string{})
	input.Filters.Page = app.readIntQueryParam(qs, "page", 1, v)
	input.Filters.PageSize = app.readIntQueryParam(qs, "page_size", 20, v)
	input.Filters.Sort = app.readStringQueryParam(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "avg_cost", "-id", "-name", "-avg_cost"}
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	breeds, metadata, err := app.models.Breeds.GetAll(input.Name, input.Countries, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if err = app.writeJSON(w, http.StatusOK, envelope{"breeds": breeds, "metadata": metadata}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
