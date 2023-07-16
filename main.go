package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sourcegraph/jsonrpc2"
)

type LSPServer struct {
	// The symmetric connection
	conn jsonrpc2.Conn

	// Check if the connection is available
	connMutex sync.Mutex

	// shutdown
	shutdown bool
}

func NewLSPServer() *LSPServer {
	return &LSPServer{}
}

func (s *LSPServer) Initialize(ctx context.Context) error {
	// Initialize here if needed
	return nil
}

func (s *LSPServer) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {
	log.Println("Handling request...")
	log.Println(req.Method)
	// Handle something
	return nil, nil
}

func (s *LSPServer) Serve(ctx context.Context) {
	log.Println("Starting LSP server...")

	// Listen on TCP port 4389 on all available unicast and
	// anycast IP addresses of the local system.
	l, err := net.Listen("tcp", "localhost:4389")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a new goroutine.
		go func(c net.Conn) {
			remote := c.RemoteAddr().String()
			log.Println("New connection: " + remote)
			// Create a new jsonrpc2 stream server
			handler := jsonrpc2.HandlerWithError(s.Handle)
			<-jsonrpc2.NewConn(
				ctx,
				jsonrpc2.NewBufferedStream(c, jsonrpc2.VSCodeObjectCodec{}),
				handler).DisconnectNotify()
			c.Close()
			log.Println("Connection closed: " + remote)
		}(conn)
	}
}

func main() {
	// Create a new LSP server
	server := NewLSPServer()
	go server.Serve(context.Background()) // run Serve in a separate goroutine
	waitForTerminationSignal()
}

func waitForTerminationSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down...")
}
