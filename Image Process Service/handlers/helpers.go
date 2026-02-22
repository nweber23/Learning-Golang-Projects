package handlers

import (
	"encoding/json"
	"net/http"

	"image-process-service/middleware"
)

func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

func GetUserIDFromContext(r *http.Request) string {
	return middleware.GetUserIDFromContext(r)
}
