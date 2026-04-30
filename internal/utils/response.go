package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Data:  data,
		Error: nil,
	})
}

func Error(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Data:  nil,
		Error: err,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Data:  nil,
		Error: message,
	})
}
