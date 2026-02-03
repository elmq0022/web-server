package handler_test

import (
	"testing"

	"github.com/elmq0022/web-server/iinternal/handler"
	"github.com/go-playground/assert/v2"
)

func TestParseHTTPRequest(t *testing.T) {
	req := []byte("POST / HTTP/1.1\r\nContent-Type: text/plain\r\n\r\nBody1\r\nBody2\r\n")

	headers := make(map[string][]string)
	headers["Content-Type"] = []string{"text/plain"}

	want := handler.HTTPRequest{
		Method:  "POST",
		URL:     "/",
		Version: "HTTP/1.1",
		Headers: headers,
		Body:    []byte("Body1\r\nBody2\r\n"),
	}

	got, err := handler.ParseRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}
