package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/elmq0022/web-server/internal/requests"
	"github.com/elmq0022/web-server/internal/responses"
)

func main() {
	host := ""
	port := "8080"
	address := net.JoinHostPort(host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Printf("server listening on %s", address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()

	req := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("error reading from %s: %v", clientAddr, err)
			return
		}
		req = append(req, buf[:n]...)
		if bytes.HasSuffix(req, []byte("\r\n\r\n")) {
			break
		}
	}

	parsedRequest := requests.ParseRequest(req)
	log.Printf("%s %s from %s", parsedRequest.Method, parsedRequest.URL, clientAddr)

	var statusCode int
	var status string
	var body []byte

	// read the page from disk
	prefix := "/var/www/"
	home, _ := os.UserHomeDir()
	page := filepath.Join(home, prefix, parsedRequest.URL)

	info, err := os.Stat(page)
	if os.IsNotExist(err) {
		log.Printf("Could not read page %s: %v", page, err)
		statusCode = http.StatusNotFound
		status = fmt.Sprintf("Page %s Not Found", parsedRequest.URL)
	} else if err != nil {
		return
	} else if info.IsDir() {
		page = filepath.Join(page, "index.html")
	} else {
		// noop
	}

	data, err := os.ReadFile(page)
	if err != nil {
		log.Printf("Could not read page %s: %v", page, err)
		statusCode = http.StatusNotFound
		status = fmt.Sprintf("Page %s Not Found", parsedRequest.URL)
	} else {
		log.Printf("Returning requested page %s", parsedRequest.URL)
		statusCode = http.StatusOK
		status = "Status OK"
		body = data
	}

	// build response
	headers := make(map[string][]string)
	headers["Content-Length"] = []string{strconv.Itoa(len(body))}

	resp := &responses.HTTPResponse{
		Version:    "HTTP/1.1",
		StatusCode: statusCode,
		Status:     status,
		Headers:    headers,
		Body:       body,
	}

	// write the response
	if _, err := conn.Write(resp.CreateResponse()); err != nil {
		log.Printf("error writing response to %s: %v", clientAddr, err)
		return
	}
	log.Printf("responded %d %s to %s", resp.StatusCode, resp.Status, clientAddr)
}
