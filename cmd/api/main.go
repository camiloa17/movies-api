package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

type application struct {
	Domain string
}

func main() {
	// set application
	var app application
	// read from command line

	// connect to db

	app.Domain = "example.com"

	server := http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		Handler:     app.routes(),
		ReadTimeout: 30 * time.Second,
	}

	// start dev server
	log.Println("startin app on port", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
