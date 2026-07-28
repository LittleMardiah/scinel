package util

import (
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"github.com/stretchr/testify/assert"
)

func TestWriteResponse(t *testing.T) {
w := httptest.NewRecorder()
requestID := "req-123"
WriteResponse(w, http.StatusOK, "Success", map[string]string{"id": "1"}, requestID)
assert.Equal(t, http.StatusOK, w.Code)
var resp ResponseWrapper
json.Unmarshal(w.Body.Bytes(), &resp)
assert.Equal(t, requestID, resp.RequestID)
}

func TestWriteError(t *testing.T) {
w := httptest.NewRecorder()
requestID := "req-456"
WriteError(w, http.StatusBadRequest, "Bad Request", "Invalid", requestID)
assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGenerateRequestID(t *testing.T) {
id1 := GenerateRequestID()
id2 := GenerateRequestID()
assert.NotEmpty(t, id1)
assert.NotEqual(t, id1, id2)
}

func TestBuildResponse(t *testing.T) {
data := map[string]interface{}{"key": "value"}
requestID := "req-789"
resp := BuildResponse(http.StatusOK, "OK", data, requestID)
assert.Equal(t, http.StatusOK, resp.Code)
assert.Equal(t, data, resp.Data)
}

func TestBuildResponse_ComplexData(t *testing.T) {
data := map[string]interface{}{
"id":    123,
"name":  "Test",
"items": []string{"a", "b"},
}
requestID := "req-123"
resp := BuildResponse(200, "OK", data, requestID)
assert.Equal(t, 200, resp.Code)
assert.Equal(t, data, resp.Data)
}

func TestErrorConstants(t *testing.T) {
assert.Equal(t, "Unauthorized", ErrUnauthorized)
assert.Equal(t, "Forbidden", ErrForbidden)
}

func TestCodeConstants(t *testing.T) {
assert.Equal(t, 200, CodeSuccess)
assert.Equal(t, 400, CodeBadRequest)
}