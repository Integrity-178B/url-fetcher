package server

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse allows to render error text as response
type ErrorResponse string

// MarshalJSON implements custom marshaling of error response to JSON
func (er ErrorResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Error string `json:"error"`
	}{
		Error: string(er),
	})
}

// ResponseWriter encapsulates response writer and add possibility to write errors and any data as JSON
type ResponseWriter struct {
	http.ResponseWriter
}

// WriteJSON writes any data to response as JSON
func (w ResponseWriter) WriteJSON(v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// WriteError writes error to response as JSON
func (w ResponseWriter) WriteError(e error) error {
	return json.NewEncoder(w).Encode(ErrorResponse(e.Error()))
}
