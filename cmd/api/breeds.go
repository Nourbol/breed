package main

import (
	"fmt"
	"github.com/Nourbol/breed/internal/data"
	"github.com/Nourbol/breed/internal/validator"
	"net/http"
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
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show the details of breed %d\n", id)
}
