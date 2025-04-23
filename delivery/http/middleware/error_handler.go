package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ErrorHandler is a middleware that handles errors
func ErrorHandler(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a custom response writer to capture errors
			ew := &errorWriter{
				ResponseWriter: w,
				logger:         logger,
			}

			// Call the next handler
			next.ServeHTTP(ew, r)
		})
	}
}

// errorWriter is a custom response writer that handles errors
type errorWriter struct {
	http.ResponseWriter
	logger *log.Logger
}

// WriteHeader overrides the WriteHeader method to log error status codes
func (ew *errorWriter) WriteHeader(code int) {
	if code >= 400 {
		ew.logger.Printf("Error: %d", code)
	}
	ew.ResponseWriter.WriteHeader(code)
}

// Error writes an error response
func (ew *errorWriter) Error(err error, message string, status int) {
	ew.logger.Printf("Error: %s", err.Error())

	// Create error response
	resp := ErrorResponse{
		Error:   err.Error(),
		Message: message,
		Status:  status,
	}

	// Write response
	ew.ResponseWriter.Header().Set("Content-Type", "application/json")
	ew.ResponseWriter.WriteHeader(status)
	json.NewEncoder(ew.ResponseWriter).Encode(resp)
}