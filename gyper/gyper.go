package gyper

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"

	"github.com/Abhishekkarunakaran/gyper/internal"
	"github.com/gofrs/uuid"
)

type HandleFunc func(c Context)

type Context struct {
	conn    net.Conn
	request *Request
}

type Job struct {
	ID      uuid.UUID
	context Context
}

type Gyper struct {
	jobChannel     chan Job
	workerCount    int
	wg             sync.WaitGroup
	listener       net.Listener
	pathMethodTree *node
}

// Creates a new gyper server
func New() (g *Gyper) {
	jobCount := 100
	return &Gyper{
		jobChannel:     make(chan Job, jobCount),
		workerCount:    runtime.NumCPU(),
		pathMethodTree: getNewNode(),
	}
}

// Start the gyper server with given ipAddr and port.
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
				ID: id,
				context: Context{
					conn:   conn,
				},
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

// Stop function closes the jobChannel and the tcp server listener
func (g *Gyper) Stop() {
	close(g.jobChannel)

	if g.listener != nil {
		err := g.listener.Close()
		if err != nil {
			fmt.Printf("Error while closing listener: %v\n", err)
		}
	}
}

// Method to register a handler function to a path as GET method
func (g *Gyper) GET(path string, function HandleFunc) {
	g.add(path, function, GET)
}

// Method to register a handler function to a path as POST method
func (g *Gyper) POST(path string, function HandleFunc) {
	g.add(path, function, POST)
}

// Method to register a handler function to a path as PUT method
func (g *Gyper) PUT(path string, function HandleFunc) {
	g.add(path, function, PUT)
}

// Method to register a handler function to a path as PATCH method
func (g *Gyper) PATCH(path string, function HandleFunc) {
	g.add(path, function, PATCH)
}

// Method to register a handler function to a path as DELETE method
func (g *Gyper) DELETE(path string, function HandleFunc) {
	g.add(path, function, DELETE)
}

// worker function to execute each rerquests
func (g *Gyper) worker() {
	defer g.wg.Done()
	for job := range g.jobChannel {
		job.context.request = NewRequest(job.context.conn)
		g.executeMethod(job.context)
		_ = job.context.conn.Close()
	}
}

// Internal function that executes the handleFunc registered in the endpoint and method type.
// If there's not function registered, it executes the defaultFunc
func (g *Gyper) executeMethod(c Context) {
	if function, exists := g.getFunction(c.request.Path, c.request.Method); !exists {
		defaultFunc(c)
	} else {
		function(c)
	}
}

// Returns the registered function in a specific path and a http method from the tree
func (g *Gyper) getFunction(path string, method Method) (HandleFunc, bool) {
	pathList := internal.GetPathList(path)
	currentNode := g.pathMethodTree
	var nextNode *node
	var exists bool
	var function HandleFunc
	for i, pathPoint := range pathList {
		nextNode, exists = currentNode.pathPoints[pathPoint]
		if !exists {
			return nil, false
		}
		currentNode = nextNode

		if len(pathList) == i+1 {
			function, exists = currentNode.methods[method]
		}
	}
	if !exists {
		return nil, false
	}
	return function, true
}

// DefaultFunc that executes whether the desired function is nott present in the tree.
func defaultFunc(c Context) {
	responseMap := map[string]string{
		"message": "Method Not Found",
	}
	response, _ := json.Marshal(responseMap)
	switch c.request.Protocol {
	case HTTP1:
		c.conn.Write([]byte("HTTP/1.1 404 Method Not Found\r\n"))
	case HTTP2:
		c.conn.Write([]byte("HTTP/2 404 Method Not Found\r\n"))
	}
	c.conn.Write([]byte("Content-Length: " + fmt.Sprint(len(response)) + "\r\n"))
	c.conn.Write([]byte("Content-Type: application/json\r\n\r\n"))
	c.conn.Write(response)
}

// Regiters the given function to the given path as the given http method.
// It is saved in the tree.
func (g *Gyper) add(path string, function HandleFunc, method Method) {
	if err := internal.ValidatePath(path); err != nil {
		log.Fatal(err.Error())
		return
	}
	if path == "" {
		path = "/"
	}

	pathList := internal.GetPathList(path)
	currentNode := g.pathMethodTree
	var nextNode *node
	var exists bool
	for i, pathPoint := range pathList {
		nextNode, exists = currentNode.pathPoints[pathPoint]
		if !exists {
			nextNode = getNewNode()
			currentNode.pathPoints[pathPoint] = nextNode
		}
		currentNode = nextNode

		if len(pathList) == i+1 {
			if currentNode.methods == nil {
				currentNode.methods = make(map[Method]HandleFunc)
			}
			currentNode.methods[method] = function
		}
	}
}

// returns new node
func getNewNode() *node {
	return &node{
		pathPoints: make(map[string]*node),
	}
}
