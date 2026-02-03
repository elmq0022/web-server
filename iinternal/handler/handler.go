package handler

import (
	"bytes"
	"strings"
)

type HTTPRequest struct {
	Method  string
	URL     string
	Version string
	Headers map[string][]string
	Body    []byte
}

func ParseRequest(req []byte) (HTTPRequest, error) {

	lines := bytes.SplitN(req, []byte("\r\n\r\n"), 2)
	headers := bytes.Split(lines[0], []byte("\r\n"))
	body := lines[1]
	l := bytes.Split(headers[0], []byte(" "))

	h := make(map[string][]string)
	for _, header := range headers[1:] {
		parts := bytes.Split(header, []byte(":"))
		k := strings.Trim(string(parts[0]), " ")
		v := strings.Trim(string(parts[1]), " ")

		if _, ok := h[k]; !ok {
			h[k] = []string{}
		}
		h[k] = append(h[k], v)
	}

	return HTTPRequest{
		Method:  string(l[0]),
		URL:     string(l[1]),
		Version: string(l[2]),
		Headers: h,
		Body:    body,
	}, nil
}

func CreateResponse(method, url string) ([]byte, error) {
	return []byte{}, nil
}
