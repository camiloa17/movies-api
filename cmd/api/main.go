package main

import (
	"context"
	"fmt"
	"log"
	"movies-backend/internal/driver"
	"movies-backend/internal/repository"
	"movies-backend/internal/repository/dbrepo"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
	Domain string
	DBRepo repository.Repository
	auth   Auth
}

func main() {
	// set application
	var app application
	// read from command line
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn, ok := os.LookupEnv("DSN")
	if !ok {
		log.Fatal("No DSN provided")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtIssuer := os.Getenv("JWT_ISSUER")
	jwtAudience := os.Getenv("JWT_AUDIENCe")
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	domain := os.Getenv("DOMAIN")
	app.Domain = domain

	app.auth = Auth{
		Issuer:        jwtIssuer,
		Audience:      jwtAudience,
		Secret:        jwtSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "Host-refresh_token",
		CookieDomain:  cookieDomain,
	}

	pgConn, err := driver.NewPgStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer pgConn.DB.Close(context.Background())

	dbRepo := dbrepo.NewStorageRepo(pgConn.DB)

	app.DBRepo = dbRepo

	server := http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		Handler:     app.routes(),
		ReadTimeout: 30 * time.Second,
	}

	// start dev server
	log.Println("startin app on port", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
