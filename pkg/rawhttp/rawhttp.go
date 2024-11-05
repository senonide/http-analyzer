package rawhttp

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/senonide/http-analyzer/pkg/ja3"
)

type RequestData struct {
	ClientIp string   `json:"client_ip"`
	JA3      string   `json:"ja3"`
	Method   string   `json:"method"`
	Headers  []string `json:"headers"`
	Body     string   `json:"body,omitempty"`
}

func (r *RequestData) HandleTcpConnection(conn net.Conn) {
	defer conn.Close()
	r.handleHTTP(conn)
}

func (r *RequestData) handleHTTP(conn net.Conn) {
	clientIp := conn.RemoteAddr().String()
	r.ClientIp = clientIp
	log.Printf("new connection from %s", clientIp)

	method, headers, body, err := readHTTPRequest(conn)
	if err != nil {
		log.Printf("error reading request: %v", err)
		return
	}

	r.Method = method
	r.Headers = headers
	r.Body = string(body)

	err = r.sendJSONResponse(conn)
	if err != nil {
		log.Printf("error sending response: %v", err)
		return
	}
}

func readHTTPRequest(conn net.Conn) (method string, headers []string, body []byte, err error) {
	reader := bufio.NewReader(conn)
	var request []byte

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", nil, nil, err
		}

		request = append(request, b)

		// Check if we have reached the end of the HTTP headers ("\r\n\r\n")
		if len(request) >= 4 && string(request[len(request)-4:]) == "\r\n\r\n" {
			break
		}
	}

	requestStr := string(request)
	lines := strings.Split(requestStr, "\r\n")

	// The first line is the request line (e.g., "GET / HTTP/1.1"), followed by headers
	for i, line := range lines {
		if i == 0 {
			method = line
			continue
		}

		// Break when encountering the blank line after the headers
		if line == "" {
			break
		}
		headers = append(headers, line)
	}

	// Check if there is a Content-Length header to read the body
	contentLength := 0
	for _, header := range headers {
		if strings.HasPrefix(strings.ToLower(header), "content-length:") {
			_, err = fmt.Sscanf(header, "Content-Length: %d", &contentLength)
			if err != nil {
				return "", nil, nil, fmt.Errorf("parsing Content-Length: %w", err)
			}
			break
		}
	}

	// If there's a body, read the specified Content-Length
	if contentLength > 0 {
		body = make([]byte, contentLength)
		_, err := io.ReadFull(reader, body)
		if err != nil {
			return "", nil, nil, err
		}
	}

	return method, headers, body, nil
}

func (r *RequestData) sendJSONResponse(conn net.Conn) error {
	jsonData, err := json.Marshal(*r)
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	response := buildHTTPResponse(len(jsonData), jsonData)

	writer := bufio.NewWriter(conn)
	if _, err := writer.WriteString(response); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flushing writer: %w", err)
	}

	return nil
}

func buildHTTPResponse(contentLength int, jsonData []byte) string {
	return "HTTP/1.1 200 OK\r\n" +
		"Content-Type: application/json\r\n" +
		"Connection: close\r\n" +
		"Content-Length: " + strconv.Itoa(contentLength) + "\r\n\r\n" +
		string(jsonData)
}

func LoadTLSConfig(cert, key string, ja3hash *string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to load key pair: %w", err)
	}

	config := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{certificate},
		GetConfigForClient: func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
			// Generate JA3 fingerprint from the ClientHello info
			*ja3hash = ja3.GenerateJA3(hello)
			return nil, nil //nolint:nilnil
		},
	}

	return config, nil
}
