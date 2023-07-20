package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"lsp.com/server/handler"
	"lsp.com/server/types"

	"github.com/sourcegraph/jsonrpc2"
)

var consoles []types.Console

type LSPServer struct {
	// The symmetric connection
	conn jsonrpc2.Conn

	shutdown bool
}

func NewLSPServer() *LSPServer {
	return &LSPServer{shutdown: false}
}

func (s *LSPServer) Run(ctx context.Context) {
	// In this case, will expose to a port. Maybe for the real
	// version it uses stdIn or similar, need to figure it out
	log.Println("Starting LSP server...")

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

func (s *LSPServer) Stop(l net.Listener) {
	l.Close()
}

func (s *LSPServer) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (result interface{}, err error) {

	var res json.RawMessage
	switch req.Method {
	case "initialize":
		log.Println("Handling initialize request...")
		output, err := handler.SendCapabilities()

		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = output

	case "shutdown":
		log.Println("Handling shutdown request...")
		s.shutdown = true
		return nil, nil

	case "exit":
		log.Println("Handling exit request...")
		if s.shutdown {
			defer conn.Close()
			return nil, nil
		} else {
			return nil, fmt.Errorf("Server not shutdown")
		}

	case "textDocument/didOpen":

		log.Println("Handling didOpen request...")
		var console types.Console

		err := json.Unmarshal(*req.Params, &console)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		consoles = append(consoles, console)
		fmt.Println(consoles)
		return nil, nil

	case "textDocument/didChange":

		log.Println("Handling didChange request...")
		var console types.Console
		err := json.Unmarshal(*req.Params, &console)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		for i := range consoles {
			if consoles[i].IP == console.IP {
				consoles[i].Content = console.Content
				fmt.Println(consoles)
				return nil, nil
			}
		}
		return nil, fmt.Errorf("IP not found")

	case "textDocument/didClose":

		log.Println("Handling didClose request...")
		var console types.Console

		err := json.Unmarshal(*req.Params, &console)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// remove from consoles the console with the given IP
		for i := range consoles {
			if consoles[i].IP == console.IP {
				consoles = append(consoles[:i], consoles[i+1:]...)
				fmt.Println(consoles)
				return nil, nil
			}
		}
		return nil, fmt.Errorf("IP not found")

	case "textDocument/diagnostics":

		log.Println("Handling diagnostics request...")
		output, err := handler.SendDiagnostics(consoles)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = output

	case "textDocument/completion":

		log.Println("Handling completion request...")
		var parameters types.CompletionRequest
		err := json.Unmarshal(*req.Params, &parameters)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		output, err := handler.SendCompletion(parameters.Console, parameters.Pointer)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		res = output

	default:
		log.Println("Unknown request...")
		return jsonrpc2.Response{ID: req.ID, Result: nil, Error: &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: "Method not found",
			Data:    nil,
		}, Meta: nil}, nil
	}

	result = jsonrpc2.Response{ID: req.ID, Result: &res, Error: nil, Meta: nil}
	return result, nil
}
