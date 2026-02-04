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

	pr := handler.ParseRequest(req)
	log.Printf("%s %s from %s", pr.Method, pr.URL, clientAddr)

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

	if _, err := conn.Write(resp.CreateResponse()); err != nil {
		log.Printf("error writing response to %s: %v", clientAddr, err)
		return
	}
	log.Printf("responded %d %s to %s", resp.StatusCode, resp.Status, clientAddr)
}
