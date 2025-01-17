package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// Verifies that the request Content-Type header is "application/json"
func CheckContentType(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		return errors.New("Content-Type must be 'application/json'")
	}
	return nil
}

// Logs the error and sends a JSON-encoded error response with the specified status code
func SendError(w http.ResponseWriter, status int, err error) {
	slog.Error(err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	err1 := encoder.Encode(map[string]string{"error": err.Error()})
	if err1 != nil {
		slog.Error("Text", slog.String("ERROR", err1.Error()))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
	}
}

// Logs a success message and sends a JSON-encoded success response with the specified status code
func SendSucces(w http.ResponseWriter, status int, message string) {
	slog.Info(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if status == http.StatusNoContent {
		return
	}
	err := encoder.Encode(map[string]string{"successfully": message})
	if err != nil {
		slog.Error("Failed to send response", slog.String("ERROR", err.Error()))
	}
}
