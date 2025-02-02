package gyper

import (
	"encoding/json"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

type HandleFunc func()
type Job struct {
	ID      uuid.UUID
	conn    net.Conn
	Request *Request
}

type Gyper struct {
	jobChannel  chan Job
	workerCount int
	wg          sync.WaitGroup
	listener    net.Listener
	tree        node
}

func New() (g *Gyper) {
	jobCount := 100
	return &Gyper{
		jobChannel:  make(chan Job, jobCount),
		workerCount: runtime.NumCPU() - 2,
		tree:        node{},
	}
}

func (g *Gyper) Start(ipAddr string, port string) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ipAddr, port))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	g.listener = listener
	defer listener.Close()

	fmt.Printf("Server listening on port: %s\n", port)

	go func() {
		defer g.wg.Done()
		g.wg.Add(1)
		for {
			conn, err := listener.Accept()

			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
					return
				}
				fmt.Println(err)
				continue
			}
			id, err := uuid.NewV4()
			if err != nil {
				fmt.Println(err)
				_ = conn.Close()
				continue
			}
			g.jobChannel <- Job{
				ID:   id,
				conn: conn,
			}
		}
	}()

	for i := 0; i < g.workerCount; i++ {
		g.wg.Add(1)
		go g.worker()
	}

	g.wg.Wait()
	return nil
}

func (g *Gyper) worker() {
	defer g.wg.Done()
	for job := range g.jobChannel {
		job.Request = NewRequest(job.conn)
		defaultFunc(job)
		job.conn.Close()
	}
}

func defaultFunc(job Job) {
	responseMap := map[string]string{
		"timstamp": time.Now().UTC().String(),
		"message":  "Method Not Found",
	}

	response, _ := json.Marshal(responseMap)

	switch job.Request.Protocol {
	case HTTP1:
		job.conn.Write([]byte("HTTP/1.1 404 Method Not Found\r\n"))
	case HTTP2:
		job.conn.Write([]byte("HTTP/2 404 Method Not Found\r\n"))
	}

	job.conn.Write([]byte("Content-Length: " + fmt.Sprint(len(response)) + "\r\n"))
	job.conn.Write([]byte("Content-Type: application/json\r\n\r\n"))
	job.conn.Write(response)
}

func (g *Gyper) Stop() {
	close(g.jobChannel)

	if g.listener != nil {
		err := g.listener.Close()
		if err != nil {
			fmt.Printf("Error while closing listener: %v\n", err)
		}
	}
	fmt.Printf("Number of goroutines:%v\n", runtime.NumGoroutine())
}
