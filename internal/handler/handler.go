package handler

import (
	"bytes"
	"maps"
	"slices"
	"strconv"
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
		parts := bytes.SplitN(header, []byte(":"), 2)
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

type HTTPResponse struct {
	Version    string
	StatusCode int
	Status     string
	Headers    map[string][]string
	Body       []byte
}

// add or update the Content-Length Header

// going to buffer the whole thing in memory to start
// caller is expected to set the body's length in the header
// Content-Length: len(r.Body)
// should really write to the stream ...
func (r *HTTPResponse) CreateResponse() []byte {
	// add start line
	resp := make([]byte, 0)
	resp = append(resp, []byte(r.Version)...)
	resp = append(resp, ' ')
	resp = append(resp, []byte(strconv.Itoa(r.StatusCode))...)
	resp = append(resp, ' ')
	resp = append(resp, []byte(r.Status)...)
	resp = append(resp, []byte("\r\n")...)

	// add headers
	// consider sorting the keys for consistent returns
	sortedKeys := slices.Sorted(maps.Keys(r.Headers))
	for _, key := range sortedKeys {
		vals, _ := r.Headers[key]
		for _, val := range vals {
			resp = append(resp, []byte(key)...)
			resp = append(resp, []byte(": ")...)
			resp = append(resp, []byte(val)...)
			resp = append(resp, []byte("\r\n")...)
		}
	}

	// append empty line
	resp = append(resp, []byte("\r\n")...)

	// append body
	resp = append(resp, r.Body...)

	return resp
}
