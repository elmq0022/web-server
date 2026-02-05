package main

import (
	"log"
	"net"
	"net/http"
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

	req, err := requests.ReadRequest(conn)
	if err != nil {
		return
	}

	parsedRequest := requests.ParseRequest(req)
	log.Printf("%s %s from %s", parsedRequest.Method, parsedRequest.URL, clientAddr)

	var resp *responses.HTTPResponse
	headers := make(map[string][]string)
	headers["Connection"] = []string{"close"}

	if parsedRequest.Method != "GET" {
		headers["Content-Length"] = []string{"0"}
		resp = &responses.HTTPResponse{
			Version:    "HTTP/1.1",
			StatusCode: http.StatusMethodNotAllowed,
			Status:     "Method Not Allowed",
			Headers:    headers,
			Body:       nil,
		}
	} else {
		statusCode, status, body := responses.LoadFileFromURL(parsedRequest.URL, headers)
		headers["Content-Length"] = []string{strconv.Itoa(len(body))}
		resp = &responses.HTTPResponse{
			Version:    "HTTP/1.1",
			StatusCode: statusCode,
			Status:     status,
			Headers:    headers,
			Body:       body,
		}
	}

	if _, err := conn.Write(resp.CreateResponse()); err != nil {
		log.Printf("error writing response to %s: %v", clientAddr, err)
		return
	}
	log.Printf("responded %d %s to %s", resp.StatusCode, resp.Status, clientAddr)
}
