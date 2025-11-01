package controllers

import (
	"encoding/json"
	"go-native-webserver/internal/apperror"
	"net/http"
)

func ResponseError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(apperror.APIError); ok {
		http.Error(w, apiErr.Message, apiErr.Code)
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ResponseSuccessJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// Encode data to JSON and write to response
	// Ignoring error handling for brevity; in production, handle errors appropriately
	json.NewEncoder(w).Encode(data)
}
