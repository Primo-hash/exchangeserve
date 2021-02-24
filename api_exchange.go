package exchangeserve

import (
	"net/http"
)

// HandlerHistory main handler for route related to `/student` requests
func HandlerHistory() func (http.ResponseWriter, *http.Request) {

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

// handleHistoryGet utility function, package level, to handle GET request to history route
func handleHistoryGet(w http.ResponseWriter, r *http.Request) {
	http.Header.Add(w.Header(), "content-type", "application/json")
	/*
	parts := strings.Split(r.URL.Path, "/")
	// error handling
	if len(parts) != 3 || parts[1] != "student" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	// handle the request /student/  which will return ALL students as array of JSON objects
	if parts[2] == "" {
		replyWithAllStudents(w, db)
	} else {
		replyWithStudent(w, db, parts[2])
	}
	*/
}