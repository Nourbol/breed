package main

import (
	"fmt"
	"github.com/Nourbol/breed/internal/data"
	"github.com/Nourbol/breed/internal/validator"
	"net/http"
	"time"
)

func (app *application) createBreedHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		AvgCost     data.Cost `json:"avg_cost"`
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
	}

	v := validator.New()
	if data.ValidateBreed(v, breed); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showBreedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	breed := data.Breed{
		ID:          id,
		CreatedAt:   time.Now(),
		Name:        "Golden Retriever",
		Description: "The Golden Retriever is a Scottish breed of retriever dog of medium size.",
		AvgCost:     2250,
		Version:     1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"breed": breed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
