package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/sourcegraph/jsonrpc2"
	"lsp.com/server/types"
)

func decodeResponseMessage(message *json.RawMessage, m map[string]interface{}) {

	err := json.Unmarshal(*message, &m)
	if err != nil {
		log.Println(err)
		return
	}

	for key, value := range m {
		fmt.Println(key + " : ")
		fmt.Println(value)
		fmt.Println("==")
	}
}

func main() {
	// connect to localhost
	conn, err := net.Dial("tcp", "localhost:4389")

	if err != nil {
		log.Fatal(err)
	}

	// create a new client
	client := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}), nil)

	// send a request
	var reply *json.RawMessage
	if err := client.Call(context.Background(), "initialize", "", &reply); err != nil {
		log.Fatal(err)
	}
	var m = map[string]interface{}{}
	decodeResponseMessage(reply, m)

	console1 := types.Console{IP: 1, Content: "aws configure"}
	console2 := types.Console{IP: 2, Content: "aws lambda list-functions"}
	console3 := types.Console{IP: 1, Content: "aws lambda list-functions --region us-east-1"}

	if err := client.Notify(context.Background(), "textDocument/didOpen", console1); err != nil {
		log.Fatal(err)
	}

	if err := client.Notify(context.Background(), "textDocument/didOpen", console2); err != nil {
		log.Fatal(err)
	}

	if err := client.Notify(context.Background(), "textDocument/didClose", console2); err != nil {
		log.Fatal(err)
	}

	if err := client.Notify(context.Background(), "textDocument/didChange", console3); err != nil {
		log.Fatal(err)
	}

	if err := client.Call(context.Background(), "textDocument/diagnostics", "", &reply); err != nil {
		log.Fatal(err)
	}

	decodeResponseMessage(reply, m)

	if err := client.Call(context.Background(), "NotExistingMethod", "", &reply); err != nil {
		log.Fatal(err)
	}

	decodeResponseMessage(reply, m)

	if err := client.Call(context.Background(), "textDocument/completion", types.CompletionRequest{
		Console: console3,
		Pointer: 0,
	}, &reply); err != nil {
		log.Fatal(err)
	}

	decodeResponseMessage(reply, m)

	if err := client.Call(context.Background(), "shutdown", "", &reply); err != nil {
		log.Fatal(err)
	}

	if err := client.Notify(context.Background(), "exit", nil); err != nil {
		log.Fatal(err)
	}

	// close the connection
	client.Close()
}
