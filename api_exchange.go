package exchange

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// p is a shortened function for extracting URL parameters
func p(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// getLimit converts string number into int and returns int, for limiting option
func getLimit(s string) int {
	// convert string number to an int and handle error for non digit characters
	if n, err := strconv.Atoi(s); err == nil {
		return n
	} else {
		// return '20' to fixed limit
		return 20
	}
}

//
func resWithData(w io.Writer, response map[string]interface{}) {
	// handle JSON objects
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}

// handleHistoryGet utility function, package level, to handle GET request to history route
func handleHistoryGet(w http.ResponseWriter, r *http.Request) {
	// Set response to be of JSON type
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 6 || parts[3] != "exchangehistory" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
	// extract URL parameters
	countryName := p(r, "country_name")
	beginDate := p(r, "b_year") + "-" + p(r, "b_month") + "-" + p(r, "b_day")
	endDate := p(r, "e_year") + "-" + p(r, "e_month") + "-" + p(r, "e_day")
	// Request currency code for country
	currencyCode, err := GetCurrency(countryName)
	if err != nil { // Error handling bad request parameter for countryName
		// In case of no server response, reply with 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}

	// Request currency history based on date period and currency code
	result, err := GetExchangeData(beginDate, endDate, currencyCode, "") // last parameter empty because not part of request
	if err != nil {                                                 // Error handling bad history request and json decoding
		// In case of no server response, reply with 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// Error could also be a 400 or failure in decoding, but we print that only internally
		fmt.Println("HTTP/JSON status: " + err.Error())
	}

	// Send result for processing
	resWithData(w, result)
}

// handleBorderGet utility function, package level, to handle GET request to border route
func handleBorderGet(w http.ResponseWriter, r *http.Request) {
	// Set response to be of JSON type
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 5 || parts[3] != "exchangeborder" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
	// extract URL parameters
	countryName := p(r, "country_name")
	// Extract optional 'limit' parameter
	number := r.URL.Query().Get("limit")
	// initiate fixed limit with most possible neighbouring countries
	limit := 20
	// r.URL.Query()["limit"] returns an array of items, so we need to choose the first item
	if len(number) > 0 { // checks if limit parameter has a value
		limit = getLimit(number)
	}

	// Request currencyBase for countryName
	currencyBase, err := GetCurrency(countryName)
	if err != nil { // Error handling bad request parameter for countryName
		// In case of no server response, reply with 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}

	currencyCode, err := GetNeighbour(countryName, limit)
	if err != nil { // Error handling bad request parameter for countryName
		// In case of no server response, reply with 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}

	// Request currency history based on date period and currency code
	result, err := GetExchangeData("", "", currencyCode, currencyBase)
	if err != nil {                                                 // Error handling bad history request and json decoding
		// In case of no server response, reply with 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// Error could also be a 400 or failure in decoding, but we print that only internally
		fmt.Println("HTTP/JSON status: " + err.Error())
	}

	// Send result for processing
	resWithData(w, result)
}

// HandlerHistory main handler for route related to `/exchangehistory` requests
func HandlerHistory() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleHistoryGet(w, r)
		case http.MethodPost:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodPut:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodDelete:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}
}

// HandlerBorder main handler for route related to `/exchangeborder` requests
func HandlerBorder() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleBorderGet(w, r)
		case http.MethodPost:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodPut:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodDelete:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}
}
