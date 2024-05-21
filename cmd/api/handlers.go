package main

import (
	"encoding/json"
	"fmt"
	"movies-backend/internal/models"
	"net/http"
	"time"
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

	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("oops something went wrong with jsoning the payload"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie

	releaseDate, err := time.Parse("2006-01-02", "1986-03-07")
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("oops something went wrong with parsing the relase date"))
		return
	}
	highLander := models.Movie{
		ID:          1,
		Title:       "High Lander",
		ReleaseDate: releaseDate,
		MPAARating:  "R",
		RunTime:     116,
		Image:       "no-image",
		Description: "Just a movie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	releaseDate, err = time.Parse("2006-01-02", "1991-06-12")
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("oops something went wrong with parsing the relase date"))
		return
	}
	rotla := models.Movie{
		ID:          2,
		Title:       "Raiders of the old Arch",
		ReleaseDate: releaseDate,
		MPAARating:  "PG-13",
		RunTime:     115,
		Description: "Another movie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	movies = append(movies, highLander, rotla)

	out, err := json.Marshal(movies)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("oops something went wrong with jsoning the payload"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
