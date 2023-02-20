package main

import (
	"fmt"
	"net/http"
)

func (app *application) createBreedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new breed")
}

func (app *application) showBreedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show the details of breed %d\n", id)
}
