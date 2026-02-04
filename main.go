package main

import (
	"bytes"
	"log"
	"net"
	"strconv"

	"github.com/elmq0022/web-server/internal/handler"
)

func main() {
	host := ""
	port := "8080"
	address := net.JoinHostPort(host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("bad: %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	req := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		req = append(req, buf[:n]...)
		if bytes.HasSuffix(req, []byte("\r\n\r\n")) {
			break
		}
	}

	// TODO: do something with the request
	pr := handler.ParseRequest(req)

	body := []byte(pr.URL)
	headers := make(map[string][]string)
	headers["Content-Length"] = []string{strconv.Itoa(len(body))}
	resp := &handler.HTTPResponse{
		Version:    "HTTP/1.1",
		StatusCode: 200,
		Status:     "OK",
		Headers:    headers,
		Body:       body,
	}
	conn.Write(resp.CreateResponse())
}
