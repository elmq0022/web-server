package requests

import (
	"bytes"
	"log"
	"net"
)

func ReadRequest(conn net.Conn) ([]byte, error) {
	clientAddr := conn.RemoteAddr().String()
	req := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("error reading from %s: %v", clientAddr, err)
			return nil, err
		}
		req = append(req, buf[:n]...)
		if bytes.HasSuffix(req, []byte("\r\n\r\n")) {
			break
		}
	}
	return req, nil
}
