package main

import (
	"exchange"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

/*
	const definitions for chi routing regex and Query string parameters
 */
const (
	// Chi regex parameters
	COUNTRY = "{country_name:[a-z]+}"			  // Country name
	BY = "{b_year:\\d\\d\\d\\d}"		   		  // Begin year
	BM = "{b_month:\\d\\d}"		   	  			  // Begin month
	BD = "{b_day:\\d\\d}"		   		          // Begin day
	EY = "{e_year:\\d\\d\\d\\d}"		   		  // End year
	EM = "{e_month:\\d\\d}"		   	              // End month
	ED = "{e_day:\\d\\d}"		   		          // End day
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
	r.Get("/", handlerHello)
	r.Get("/exchange/v1/exchangehistory/"+COUNTRY+"/"+BY+"-"+BM+"-"+BD+"-"+EY+"-"+EM+"-"+ED, exchange.HandlerHistory())
	r.Get("/exchange/v1/exchangeborder/"+COUNTRY, exchange.HandlerBorder())
	//r.Get("/exchange/v1/diag/", http.HandlerFunc())

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func handlerHello(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		status := http.StatusBadRequest
		http.Error(w, "Expecting format .../firstname/lastname", status)
		return
	}
	name := parts[2]
	_, err := fmt.Fprintln(w, parts)
	if err != nil {
		// TODO must handle the error!
	}
	_, err = fmt.Fprintf(w, "Hello %s %s!\n", name, parts[3])
	if err != nil {
		// TODO must handle the error!
	}
}