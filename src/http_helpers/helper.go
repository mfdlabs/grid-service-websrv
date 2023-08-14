package httphelpers

// This file is for simple response helpers.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// WriteJSON writes a JSON response.
func WriteJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

// WriteError writes an error response.
func WriteError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
	})
}

// WriteRobloxJSONError writes a Roblox JSON error response.
func WriteRobloxJSONError(w http.ResponseWriter, str string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"code":    1,
			"message": str,
		},
	})
}

// WriteRobloxJSONErr writes a Roblox JSON error response.
func WriteRobloxJSONErr(w http.ResponseWriter, err error) {
	WriteRobloxJSONError(w, err.Error())
}

// ReadJSON reads a JSON request.
func ReadJSON(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteRobloxJSONError(w, fmt.Sprintf("Read input: %v", err))
		return false
	}
	return true
}

// ParseInt64FromQuery parses an int64 from a query.
func ParseInt64FromQuery(r *http.Request, key string) (int64, error) {
	return ParseInt64(r.URL.Query().Get(key))
}

// ParseInt64 parses an int64.
func ParseInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}
