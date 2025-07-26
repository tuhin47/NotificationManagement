package types

import "time"

// ErrorResponse represents a standardized error response structure
type ErrorResponse struct {
	Message    string    `json:"message"`
	Error      string    `json:"error"`
	StatusCode int       `json:"status_code"`
	Timestamp  time.Time `json:"timestamp"`
	ErrorCode  string    `json:"error_code"`
}
