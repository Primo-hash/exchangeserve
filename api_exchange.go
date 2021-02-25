package exchange

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strings"
)

func resWithData(w io.Writer, response map[string]interface{}) {
	// handle JSON objects
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}


// handleHistoryGet utility function, package level, to handle GET request to history route
func HandleHistoryGet(w http.ResponseWriter, r *http.Request) {
	// Set response to be of JSON type
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 6 || parts[3] != "exchangehistory" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
	// shortened function for extracting URL parameters
	p := func(r *http.Request, key string) string {
		return chi.URLParam(r, key)
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
	result, err := GetHistory(beginDate, endDate, currencyCode, "")	// last parameter empty because not part of request
	if err != nil { // Error handling bad history request and json decoding
		// In case of no server response, reply with 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// Error could also be a 400 or failure in decoding, but we print that only internally
		fmt.Println("HTTP/JSON status: " + err.Error())
	}

	// Send result for processing
	resWithData(w, result)
}


// HandlerHistory main handler for route related to `/exchangehistory` requests
func HandlerHistory() func (http.ResponseWriter, *http.Request) {

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			HandleHistoryGet(w, r)
		case http.MethodPost:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodPut:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodDelete:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
	}

	return httpHandler
}