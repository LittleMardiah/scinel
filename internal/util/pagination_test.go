package util

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedLimit  int
		expectedOffset int
	}{
		{"defaults", map[string]string{}, 20, 0},
		{"limit only", map[string]string{"limit": "50"}, 50, 0},
		{"offset only", map[string]string{"offset": "10"}, 20, 10},
		{"both", map[string]string{"limit": "30", "offset": "5"}, 30, 5},
		{"invalid limit", map[string]string{"limit": "abc"}, 20, 0},
		{"invalid offset", map[string]string{"offset": "xyz"}, 20, 0},
		{"negative offset", map[string]string{"offset": "-5"}, 20, 0},
		{"limit zero", map[string]string{"limit": "0"}, 20, 0},
		{"limit max", map[string]string{"limit": "100"}, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{URL: &url.URL{RawQuery: ""}}
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			limit, offset := ParsePagination(req)
			assert.Equal(t, tt.expectedLimit, limit)
			assert.Equal(t, tt.expectedOffset, offset)
		})
	}
}

func TestParsePagination_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		query    map[string]string
		expectedLimit, expectedOffset int
	}{
		{"limit_negative", map[string]string{"limit": "-5"}, 20, 0},
		{"offset_negative", map[string]string{"offset": "-5"}, 20, 0},
		{"limit_string", map[string]string{"limit": "abc"}, 20, 0},
		{"offset_string", map[string]string{"offset": "xyz"}, 20, 0},
		{"limit_large", map[string]string{"limit": "1000"}, 1000, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{URL: &url.URL{RawQuery: ""}}
			q := req.URL.Query()
			for k, v := range tt.query {
				q.Set(k, v)
			}
			req.URL.RawQuery = q.Encode()
			limit, offset := ParsePagination(req)
			assert.Equal(t, tt.expectedLimit, limit)
			assert.Equal(t, tt.expectedOffset, offset)
		})
	}
}
