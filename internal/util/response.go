package util

import (
"encoding/json"
"net/http"
"time"

"github.com/google/uuid"
)

// ResponseWrapper is the standard response format (TDD Section 5.1)
type ResponseWrapper struct {
Code      int         `json:"code"`
Message   string      `json:"message"`
Data      interface{} `json:"data,omitempty"`
Details   string      `json:"details,omitempty"`
Timestamp string      `json:"timestamp"`
RequestID string      `json:"request_id"`
}

// ErrorCodes (TDD Section 5.1.2)
const (
CodeSuccess           = 200
CodeCreated           = 201
CodeBadRequest        = 400
CodeUnauthorized      = 401
CodeForbidden         = 403
CodeNotFound          = 404
CodeConflict          = 409
CodeInternalError     = 500
CodeServiceUnavailable = 503
)

// Error messages
var (
ErrUnauthorized   = "Unauthorized"
ErrForbidden      = "Forbidden"
ErrNotFound       = "Not Found"
ErrBadRequest     = "Bad Request"
ErrConflict       = "Conflict"
ErrInternalError  = "Internal Server Error"
)

// WriteResponse writes a successful response
func WriteResponse(w http.ResponseWriter, code int, message string, data interface{}, requestID string) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(code)

response := ResponseWrapper{
Code:      code,
Message:   message,
Data:      data,
Timestamp: time.Now().UTC().Format(time.RFC3339),
RequestID: requestID,
}

json.NewEncoder(w).Encode(response)
}

// WriteError writes an error response
func WriteError(w http.ResponseWriter, code int, message string, details string, requestID string) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(code)

response := ResponseWrapper{
Code:      code,
Message:   message,
Details:   details,
Timestamp: time.Now().UTC().Format(time.RFC3339),
RequestID: requestID,
}

json.NewEncoder(w).Encode(response)
}

// GenerateRequestID generates a unique request ID
func GenerateRequestID() string {
return uuid.New().String()
}
