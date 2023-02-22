package main

import (
	"fmt"
	"github.com/Nourbol/breed/internal/data"
	"net/http"
	"time"
)

func (app *application) createBreedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new breed")
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
