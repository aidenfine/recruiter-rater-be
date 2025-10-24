package common

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func Ok(w http.ResponseWriter, data any) {
	respond(w, http.StatusOK, APIResponse{
		Status: "OK",
		Data:   data,
	})
}

// func that returns an http response with an error
func Error(w http.ResponseWriter, code int, message string) {
	respond(w, code, APIResponse{
		Status:  "ERROR",
		Message: message,
	})
}

func DecodeJSONBody[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var body T
	err := json.NewDecoder(r.Body).Decode(&body)
	return body, err
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respond(w http.ResponseWriter, code int, body APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
