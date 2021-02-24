package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)


func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Define new router
	r := chi.NewRouter()
	// Logs the start and end of each request with the elapsed processing time
	r.Use(middleware.Logger)

	r.Get("/exchange/v1/exchangehistory/")
	r.Get("/exchange/v1/exchangeborder/")
	r.Get("/exchange/v1/diag/")

	log.Fatal(http.ListenAndServe(":"+port, r))
}