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
		Name string `json:"name"`
		Year int32  `json:"year"`
		// something
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	breed := &data.Breed{
		// data breed
	}

	v := validator.New()

	//something

	if data.ValidateBreed(v, breed); !v.valid() {
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
		Name:        "  ",
		Description: "",
		AvgCost:     ' ',
		Version:     ' ',
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"breed": breed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
