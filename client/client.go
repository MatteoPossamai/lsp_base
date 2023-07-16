package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/sourcegraph/jsonrpc2"
)

func main() {
	// connect to localhost:5000
	conn, err := net.Dial("tcp", "localhost:4389")

	if err != nil {
		log.Fatal(err)
	}

	// create a new client
	client := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), nil)

	// send a request
	var reply string
	if err := client.Call(context.Background(), "email.SendEmail", "Hello World!", &reply); err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)

	// close the connection
	client.Close()
}
