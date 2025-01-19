package gyper

import (
	"fmt"
	"net"
	"runtime"
	"sync"

	"github.com/gofrs/uuid"
)

type Job struct {
	ID   uuid.UUID
	conn net.Conn
}

type Gyper struct {
	jobChannel  chan Job
	workerCount int
	wg          sync.WaitGroup
	listener    net.Listener
}

func New() (g *Gyper) {
	jobCount := 100
	return &Gyper{
		jobChannel:  make(chan Job, jobCount),
		workerCount: runtime.NumCPU()-2,
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
		defaultFunc(job.conn, job.ID.String())
	}
}

func defaultFunc(conn net.Conn, content string) {
	defer conn.Close()

	// time.Sleep(500 * time.Millisecond)
	response := content + "\r\n"
	conn.Write([]byte("HTTP/2 200 OK\r\n"))
	conn.Write([]byte("Content-Length: " + fmt.Sprint(len(response)) + "\r\n"))
	conn.Write([]byte("Content-Type: text/plain\r\n\r\n"))
	conn.Write([]byte(response))

}

func (g *Gyper) Stop() {
	close(g.jobChannel)

	if g.listener != nil {
		err := g.listener.Close()
		if err != nil {
			fmt.Printf("Error while closing listener: %v\n", err)
		}
	}
	fmt.Printf("Number of goroutines:%v\n",runtime.NumGoroutine())
}
