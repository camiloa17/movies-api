package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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

func (app *application) Authenticate(w http.ResponseWriter, r *http.Request) {
	// read a json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.ReadJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	// validate user against DB
	user, err := app.DBRepo.GetUserByEmail(requestPayload.Email)
	if err != nil || user == nil {
		log.Println(err)
		app.ErrorJSON(w, errors.New("password or user where incorrect, please try again"), http.StatusBadRequest)
		return
	}
	// check password hash
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Println(err)
		app.ErrorJSON(w, errors.New("password or user where incorrect, please try again"), http.StatusBadRequest)
		return
	}
	// create a jwt user.
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// generate tokens
	tokensPair, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err, http.StatusInternalServerError)
	}

	refreshCookie := app.auth.GetRefreshCookie(tokensPair.RefreshToken)

	http.SetCookie(w, refreshCookie)

	app.WriteJSON(w, http.StatusAccepted, tokensPair)
}

func (app *application) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var cookie *http.Cookie
	for _, c := range r.Cookies() {
		if c.Name == app.auth.CookieName {
			cookie = c
			break
		}
	}

	if cookie == nil {
		app.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	claims := &Claims{}
	refreshToken := cookie.Value

	// parse to token
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		return []byte(app.auth.Secret), nil
	})
	if err != nil {
		app.ErrorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}
	// get user id from claims
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.ErrorJSON(w, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}

	user, err := app.DBRepo.GetUserByID(userID)
	if err != nil {
		app.ErrorJSON(w, errors.New("not authorized"), http.StatusUnauthorized)
		return
	}

	jwu := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	tokenPairs, err := app.auth.GenerateTokenPair(&jwu)
	if err != nil {
		app.ErrorJSON(w, errors.New("no tokens generated"), http.StatusUnauthorized)
		return
	}
	refreshCookie := app.auth.GetRefreshCookie(tokenPairs.RefreshToken)

	http.SetCookie(w, refreshCookie)
	app.WriteJSON(w, http.StatusOK, tokenPairs)
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) MovieCatalog(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DBRepo.AllMovies()
	if err != nil {
		fmt.Println(err)
		app.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = app.WriteJSON(w, http.StatusOK, movies)
}
