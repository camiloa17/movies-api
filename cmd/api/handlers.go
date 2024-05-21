package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go movies up and running",
		Version: "v1.0.0",
	}

	_ = app.WriteJSON(w, http.StatusOK, payload)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DBRepo.AllMovies()
	if err != nil {
		fmt.Println(err)
		app.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = app.WriteJSON(w, http.StatusOK, movies)
}

func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		fmt.Println(err)
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	movie, err := app.DBRepo.GetMovie(id)
	if err != nil {
		fmt.Println(err)
		app.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = app.WriteJSON(w, http.StatusOK, movie)
}
