package gyper

import (
	"fmt"
	"net"
)

type Gyper struct{}

func New() (g *Gyper) {
	return &Gyper{}
}

func (g *Gyper) Start(ipAddr string, port string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ipAddr, port))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer listener.Close()

	fmt.Printf("Server listening on port: %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go defaultFunc(conn)
	}
}

func defaultFunc(conn net.Conn) {
	defer conn.Close()
	response := "Hello, World!\r\n"
	conn.Write([]byte("HTTP/2 200 OK\r\n"))
	conn.Write([]byte("Content-Length: " + fmt.Sprint(len(response)) + "\r\n"))
	conn.Write([]byte("Content-Type: text/plain\r\n\r\n"))
	conn.Write([]byte(response))
}
