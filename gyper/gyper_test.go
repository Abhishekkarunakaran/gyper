package gyper

import (
	"bufio"
	"net"
	"sync"
	"testing"
	"time"
)

func startServer() *Gyper {
	server := New()
	go func() {
		server.Start("localhost", "8888")
	}()

	time.Sleep(1 * time.Second) // Give the server time to start
	return server
}

func simulateClient(wg *sync.WaitGroup, b *testing.B) {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		b.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		b.Fatalf("Failed to read response: %v", err)
	}

	// Check response content
	if len(response) == 0 {
		b.Fatalf("Empty response received from server")
	}
}

func BenchmarkGyperServer(b *testing.B) {
	server := startServer()
	defer server.Stop()

	var wg sync.WaitGroup

	// Simulate `b.N` concurrent clients
	b.ResetTimer() // Reset the timer to exclude server startup time
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go simulateClient(&wg, b)
	}
	wg.Wait()
}
