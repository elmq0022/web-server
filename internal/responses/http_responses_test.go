package responses_test

import (
	"testing"

	"github.com/elmq0022/web-server/internal/responses"
	"github.com/go-playground/assert/v2"
)

func TestCreateResponse(t *testing.T) {
	tests := []struct {
		name string
		resp responses.HTTPResponse
		want []byte
	}{
		{
			name: "200 OK with no headers and no body",
			resp: responses.HTTPResponse{
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
			resp: responses.HTTPResponse{
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
			resp: responses.HTTPResponse{
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
			resp: responses.HTTPResponse{
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
			resp: responses.HTTPResponse{
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
