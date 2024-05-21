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

	pgConn, err := driver.NewPgStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pgConn.DB.Close(context.Background())

	dbRepo := dbrepo.NewStorageRepo(pgConn.DB)

	app.Domain = "example.com"
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
