package main

import (
	"context"
	"fmt"
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
	fmt.Println("Handling request...")
	fmt.Println(req.ID)
	// Handle something
	return nil, nil
}

func (s *LSPServer) Serve(ctx context.Context) {
	fmt.Println("Starting LSP server...")

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
			fmt.Println("New connection: " + remote)
			// Create a new jsonrpc2 stream server
			handler := jsonrpc2.HandlerWithError(s.Handle)
			<-jsonrpc2.NewConn(
				ctx,
				jsonrpc2.NewBufferedStream(c, jsonrpc2.VSCodeObjectCodec{}),
				handler).DisconnectNotify()
			c.Close()
			fmt.Println("Connection closed: " + remote)
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

// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/sourcegraph/jsonrpc2"
// )

// type handler struct{}

// func (h *handler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
// 	fmt.Println("Handle")
// 	switch req.Method { // req.Method
// 	case "email.SendEmail":
// 		h.getStatus(ctx, conn, req)
// 	default:
// 		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{Code: 1, Message: fmt.Sprintf("Method %q not found", req.Method), Data: nil})
// 	}
// 	fmt.Println("Handle done")
// }

// func (h *handler) getStatus(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) error {
// 	fmt.Println("getStatus")
// 	id := jsonrpc2.ID{}
// 	conn.Reply(ctx, id, "OK")
// 	fmt.Println("getStatus done")
// 	return nil
// }

// type Server struct {
// 	// The connection object
// 	conn *jsonrpc2.Conn

// 	// The handler
// 	handler *handler
// }

// func NewServer() *Server {
// 	return &Server{}
// }

// func (s *Server) Start() {
// 	// Create a handler
// 	s.handler = &handler{}

// 	// Start listening for requests

// }

// func (s *Server) Serve(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
// 	s.conn = conn
// 	// when a request is received, the handler is called

// 	s.handler.Handle(ctx, conn, req)
// }

// func main() {
// 	// Create a listener for the server
// 	listener, err := net.Listen("tcp", "localhost:5000")
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}

// 	fmt.Println("Server started on localhost:5000")

// 	defer listener.Close()

// 	go func() {
// 		for {
// 			conn, err := listener.Accept()
// 			if err != nil {
// 				log.Printf("Failed to accept connection: %v", err)
// 				continue
// 			}

// 			opt := jsonrpc2.OnRecv(func(req *jsonrpc2.Request, res *jsonrpc2.Response) {
// 				fmt.Println("OnRecv")
// 			})
// 			rpcConn := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewPlainObjectStream(conn), &handler{}, opt)

// 			fmt.Println("rpcConn", rpcConn)

// 		}
// 	}()

// 	// Wait for a termination signal to gracefully shutdown the server
// 	waitForTerminationSignal()
// }
