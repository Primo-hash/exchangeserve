package exchange

import (
	"encoding/json"
	"exchange/country"
	"exchange/currency"
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var appStart time.Time

// Diagnose struct for JSON encoding
type Diagnose struct {
	Exchangeratesapi string `json:"exchangeratesapi"`
	Restcountries    string `json:"restcountries"`
	Version          string `json:"version"`
	Uptime           string `json:"uptime"`
}

/*
HandlerLostUser is a function for guiding lost souls back to the right 'relative' path
*/
func HandlerLostUser(w http.ResponseWriter, r *http.Request) {
	protocol := "http://"
	host := r.Host 					// Host URL
	pathAPI := "/exchange/v1/" 		// API path

	// API request queries
	pathHistory := "exchangehistory/norway/2020-03-01-2020-03-03"
	pathBorder := "exchangeborder/norway?limit=2"
	pathDiag := "diag/"

	line1 := "Hello! You seem lost! Let me help you!"
	line2 := "These are some examples:"
	history := protocol + host + pathAPI + pathHistory
	border := protocol + host + pathAPI + pathBorder
	diagnose := protocol + host + pathAPI + pathDiag

	// HTML form for response such that URLs are hyperlinks
	var form = `<p>`+line1+`</p>
				<p>`+line2+`</p>
			    <p><a href="`+history+`">`+history+`</a></p>
				<p><a href="`+border+`">`+border+`</a></p>
				<p><a href="`+diagnose+`">`+diagnose+`</a></p>`

	// Generate HTML template from Form
	res := template.New("table")
	res, err := res.Parse(form)
	if err != nil {
		http.Error(w, "Could not process request", http.StatusInternalServerError)
		fmt.Println("Could not parse HTML form: " + err.Error())
	}

	// Write template
	err = res.Execute(w, nil)
	if err != nil {
		http.Error(w, "Could not process request", http.StatusInternalServerError)
		fmt.Println("Could not write back response: " + err.Error())
	}
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

// HandlerDiag main handler for route related to `/diag` requests
func HandlerDiag(t time.Time) func(http.ResponseWriter, *http.Request) {
	appStart = t	// Pass application start time for multiple function access
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleDiagGet(w, r)
		case http.MethodPost:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodPut:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		case http.MethodDelete:
			http.Error(w, "Not implemented", http.StatusNotImplemented)
		}
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
	currencyCode, err := country.GetCurrency(countryName)
	if err != nil { // Error handling bad request parameter for countryName
		// In case of no server response, reply with 500
		http.Error(w, "Could not contact API server", http.StatusInternalServerError)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}

	// Request currency history based on date period and currency code
	result, err := currency.GetExchangeData(beginDate, endDate, currencyCode, "") // last parameter empty because not part of request
	if err != nil {                                                               // Error handling bad history request and json decoding
		// In case of no server response, reply with 500
		http.Error(w, "Could not contact API server", http.StatusInternalServerError)
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
	currencyBase, err := country.GetCurrency(countryName)
	if err != nil { // Error handling bad request parameter for countryName
		// In case of no server response, reply with 500
		http.Error(w, "Could not contact API server", http.StatusInternalServerError)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}

	currencyCode, err := country.GetNeighbour(countryName, limit)
	if err != nil { // Error handling bad request parameter for countryName
		// In case of no server response, reply with 500
		http.Error(w, "Could not contact API server", http.StatusInternalServerError)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}

	// Request currency history based on date period and currency code
	result, err := currency.GetExchangeData("", "", currencyCode, currencyBase)
	if err != nil {                                                 // Error handling bad history request and json decoding
		// In case of no server response, reply with 500
		http.Error(w, "Could not contact API server", http.StatusInternalServerError)
		// Error could also be a 400 or failure in decoding, but we print that only internally
		fmt.Println("HTTP/JSON status: " + err.Error())
	}

	// Send result for processing
	resWithData(w, result)
}

// handleDiagGet utility function, package level, to handle GET request to diag route
func handleDiagGet(w http.ResponseWriter, r *http.Request) {
	var diag Diagnose
	var err error
	// Set response to be of JSON type
	http.Header.Add(w.Header(), "content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 5 || parts[3] != "diag" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
	// Insert exchangerate status code
	diag.Exchangeratesapi, err = currency.HealthCheck()
	if err != nil {
		// Error could be a 400, print internally as well
		fmt.Println("HTTP status: " + err.Error())
	}
	// Insert restcountries status code
	diag.Restcountries, err = country.HealthCheck()
	if err != nil {
		// Error could be a 400, print internally as well
		fmt.Println("HTTP status: " + err.Error())
	}
	// Insert API version
	diag.Version = parts[2]
	// Insert API uptime in hr min sec
	diag.Uptime = time.Since(appStart).String()
	// Encode diagnostic report
	report, _ := json.Marshal(diag)
	if err != nil {
		// In case of no server response, reply with 500
		http.Error(w, "Could not process request", http.StatusInternalServerError)
		// Error could also be a 400 or failure in decoding, but we print that only internally
		fmt.Println("Encode: " + err.Error())
	}
	// Send status and diagnostic report
	w.WriteHeader(http.StatusOK)
	w.Write(report)		 		 // Send result for processing
}

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

// resWithData write map encoded as a JSON to http response
func resWithData(w io.Writer, response map[string]interface{}) {
	// handle JSON objects
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}