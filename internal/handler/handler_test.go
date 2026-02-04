package handler_test

import (
	"testing"

	"github.com/elmq0022/web-server/internal/handler"
	"github.com/go-playground/assert/v2"
)

func TestParseHTTPRequest(t *testing.T) {
	tests := []struct {
		name string
		req  []byte
		want handler.HTTPRequest
	}{
		{
			name: "POST with body",
			req:  []byte("POST / HTTP/1.1\r\nContent-Type: text/plain\r\n\r\nBody1\r\nBody2\r\n"),
			want: handler.HTTPRequest{
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
			want: handler.HTTPRequest{
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
			want: handler.HTTPRequest{
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
			want: handler.HTTPRequest{
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
			want: handler.HTTPRequest{
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
			got := handler.ParseRequest(tt.req)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateResponse(t *testing.T) {
	tests := []struct {
		name string
		resp handler.HTTPResponse
		want []byte
	}{
		{
			name: "200 OK with no headers and no body",
			resp: handler.HTTPResponse{
				Version:    "HTTP/1.1",
				StatusCode: 200,
				Status:     "OK",
				Headers:    map[string][]string{},
				Body:       []byte(""),
			},
			want: []byte("HTTP/1.1 200 OK\r\n\r\n"),
		},
		{
			name: "200 OK with multiple headers and body",
			resp: handler.HTTPResponse{
				Version:    "HTTP/1.1",
				StatusCode: 200,
				Status:     "OK",
				Headers:    map[string][]string{"Content-Type": {"text/html"}, "Content-Length": {"13"}, "X-Custom": {"value"}},
				Body:       []byte("<html></html>"),
			},
			want: []byte("HTTP/1.1 200 OK\r\nContent-Length: 13\r\nContent-Type: text/html\r\nX-Custom: value\r\n\r\n<html></html>"),
		},
		{
			name: "404 Not Found with single header",
			resp: handler.HTTPResponse{
				Version:    "HTTP/1.1",
				StatusCode: 404,
				Status:     "Not Found",
				Headers:    map[string][]string{"Content-Length": {"0"}},
				Body:       []byte(""),
			},
			want: []byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n"),
		},
		{
			name: "201 Created with multiple headers and JSON body",
			resp: handler.HTTPResponse{
				Version:    "HTTP/1.1",
				StatusCode: 201,
				Status:     "Created",
				Headers:    map[string][]string{"Content-Type": {"application/json"}, "Location": {"/api/users/1"}, "Cache-Control": {"no-cache"}},
				Body:       []byte("{\"id\":1}"),
			},
			want: []byte("HTTP/1.1 201 Created\r\nCache-Control: no-cache\r\nContent-Type: application/json\r\nLocation: /api/users/1\r\n\r\n{\"id\":1}"),
		},
		{
			name: "500 Internal Server Error",
			resp: handler.HTTPResponse{
				Version:    "HTTP/1.1",
				StatusCode: 500,
				Status:     "Internal Server Error",
				Headers:    map[string][]string{},
				Body:       []byte("Something went wrong"),
			},
			want: []byte("HTTP/1.1 500 Internal Server Error\r\n\r\nSomething went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.resp.CreateResponse()
			assert.Equal(t, tt.want, got)
		})
	}
}
