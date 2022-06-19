package api

import (
	"encoding/json"
	_ "github.com/gorilla/mux"
	"net/http"
)

func initHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

type Message struct {
	StatusCode int    `json:"status code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

func ResponseJSON(w http.ResponseWriter, statusCode int, isError bool, message string) {
	var msg Message
	msg.StatusCode = statusCode
	msg.Message = message
	msg.IsError = isError

	w.WriteHeader(400)
	json.NewEncoder(w).Encode(msg)
}
