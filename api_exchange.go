package exchangeserve

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

func replyWithData(w io.Writer, response map[string]interface{}) {
	// handle JSON objects
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
}


// handleHistoryGet utility function, package level, to handle GET request to history route
func HandleHistoryGet(w http.ResponseWriter, r *http.Request) {
	println("four")
	// Set response to be of JSON type
	http.Header.Add(w.Header(), "content-type", "application/json")
	/*parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 8 || parts[5] != "exchangehistory" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
	*/
	// extract URL parameters
	countryName := chi.URLParam(r, "country_name")
	beginDate := chi.URLParam(r, "begin_date")
	endDate := chi.URLParam(r, "end_date")
	// Request currency code for country
	println("five")
	currencyCode, err := GetCurrency(countryName)
	println("six")

	if err != nil { // Error handling bad request parameter for countryName
		// In case of no server response, reply with 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// Error could also be a 400, but we print that only internally
		fmt.Println("HTTP status: " + err.Error())
	}
	/*
	// Request currency history based on date period and currency code
	result, err := GetHistory(beginDate, endDate, currencyCode, "")	// last parameter empty because not part of request
	if err != nil { // Error handling bad history request and json decoding
		// In case of no server response, reply with 404
		http.Error(w, err.Error(), http.StatusNotFound)
		// Error could also be a 400 or failure in decoding, but we print that only internally
		fmt.Println("HTTP/JSON status: " + err.Error())
	}
	*/
	_, err = fmt.Fprintf(w, "%s %s / %s", currencyCode, beginDate, endDate)
	if err != nil {
		// TODO must handle the error!
	}
	/*
	// handle JSON objects
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		fmt.Println("ERROR encoding JSON", err)
	}
	*/

	// Send result for processing
	//replyWithData(w, result)
}


// HandlerHistory main handler for route related to `/exchangehistory` requests
func HandlerHistory() func (http.ResponseWriter, *http.Request) {

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		println("two")
		switch r.Method {
		case http.MethodGet:
			println("three")
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