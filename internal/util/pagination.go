package util

import (
    "net/http"
    "strconv"
)

// ParsePagination extracts limit and offset from request query parameters
// Default: limit=20, offset=0
func ParsePagination(r *http.Request) (limit, offset int) {
    limit = 20
    offset = 0
    if l := r.URL.Query().Get("limit"); l != "" {
        if v, err := strconv.Atoi(l); err == nil && v > 0 {
            limit = v
        }
    }
    if o := r.URL.Query().Get("offset"); o != "" {
        if v, err := strconv.Atoi(o); err == nil && v >= 0 {
            offset = v
        }
    }
    return
}