package gyper

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
)

type Method string

const (
	//method to get data. (read only)
	GET Method = "GET"

	//method to save/send data to the server.
	POST Method = "POST"

	//method to update an entity in the server.
	PUT Method = "PUT"

	//method to partial update entity in the server
	PATCH Method = "PATCH"

	//method to delete an entity from the server.
	DELETE Method = "DELETE"

	//method to get metadata.
	HEAD Method = "HEAD"

	//method to get information about the possible communication options.
	OPTIONS Method = "OPTIONS"

	//method is for diagnosis purposes.
	TRACE Method = "TRACE"

	//method is for making end-to-end connections between a client and a server.
	CONNECT Method = "CONNECT"
)

type Protocol string

const (
	HTTP1 Protocol = "HTTP/1.1"
	HTTP2 Protocol = "HTTP/2" 
)

var protocolMap = map[string]Protocol{
	"HTTP/1.1" : HTTP1,
	"HTTP/2" : HTTP2,
}

func getProtocol(protocol string) Protocol {
	return protocolMap[protocol]
}

var methodMap = map[string]Method{
	"GET":     GET,
	"POST":    POST,
	"PUT":     PUT,
	"PATCH":   PATCH,
	"DELETE":  DELETE,
	"HEAD":    HEAD,
	"OPTIONS": OPTIONS,
	"TRACE":   TRACE,
	"CONNECT": CONNECT,
}

func getMethod(methodName string) Method {
	return methodMap[methodName]
}

type Request struct {
	Method   Method
	Path     string
	Protocol Protocol
	Header   map[string]string
	Body     io.Reader
}

func NewRequest(conn net.Conn) *Request {
	reader := bufio.NewReader(conn)

	firstLine := true
	headers := make(map[string]string, 0)
	var request Request
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return nil
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		if firstLine {
			parts := strings.Split(line, " ")
			if len(parts) == 3 {
				request.Method = getMethod(parts[0])
				request.Path = parts[1]
				request.Protocol = getProtocol(parts[2])
			}
			firstLine = false
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	request.Header = headers
	var contentLength int
	if val, ok := request.Header["Content-Length"]; ok {
		_, _ = fmt.Sscanf(val, "%d", &contentLength)
	}

	body := make([]byte, contentLength)
	_, err := reader.Read(body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	request.Body = bytes.NewReader(body)

	return &request
}
