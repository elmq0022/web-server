package requests_test

import (
	"testing"

	"github.com/elmq0022/web-server/internal/requests"
	"github.com/go-playground/assert/v2"
)

func TestParseHTTPRequest(t *testing.T) {
	tests := []struct {
		name string
		req  []byte
		want requests.HTTPRequest
	}{
		{
			name: "POST with body",
			req:  []byte("POST / HTTP/1.1\r\nContent-Type: text/plain\r\n\r\nBody1\r\nBody2\r\n"),
			want: requests.HTTPRequest{
				Method:  "POST",
				URL:     "/",
				Version: "HTTP/1.1",
				Headers: map[string][]string{"Content-Type": {"text/plain"}},
				Body:    []byte("Body1\r\nBody2\r\n"),
			},
		},
		{
			name: "GET without body",
			req:  []byte("GET /index.html HTTP/1.1\r\nHost: example.com\r\n\r\n"),
			want: requests.HTTPRequest{
				Method:  "GET",
				URL:     "/index.html",
				Version: "HTTP/1.1",
				Headers: map[string][]string{"Host": {"example.com"}},
				Body:    []byte(""),
			},
		},
		{
			name: "GET with query parameters",
			req:  []byte("GET /search?q=golang&page=1 HTTP/1.1\r\nHost: example.com\r\nAccept: text/html\r\n\r\n"),
			want: requests.HTTPRequest{
				Method:  "GET",
				URL:     "/search?q=golang&page=1",
				Version: "HTTP/1.1",
				Headers: map[string][]string{"Host": {"example.com"}, "Accept": {"text/html"}},
				Body:    []byte(""),
			},
		},
		{
			name: "POST with JSON body",
			req:  []byte("POST /api/users HTTP/1.1\r\nContent-Type: application/json\r\nAuthorization: Bearer token123\r\n\r\n{\"name\":\"John\"}"),
			want: requests.HTTPRequest{
				Method:  "POST",
				URL:     "/api/users",
				Version: "HTTP/1.1",
				Headers: map[string][]string{"Content-Type": {"application/json"}, "Authorization": {"Bearer token123"}},
				Body:    []byte("{\"name\":\"John\"}"),
			},
		},
		{
			name: "DELETE request",
			req:  []byte("DELETE /api/users/42 HTTP/1.1\r\nHost: api.example.com\r\n\r\n"),
			want: requests.HTTPRequest{
				Method:  "DELETE",
				URL:     "/api/users/42",
				Version: "HTTP/1.1",
				Headers: map[string][]string{"Host": {"api.example.com"}},
				Body:    []byte(""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := requests.ParseRequest(tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}
