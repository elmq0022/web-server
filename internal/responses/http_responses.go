package responses

import (
	"maps"
	"slices"
	"strconv"
)

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
