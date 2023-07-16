package handler

import (
	"context"
	"fmt"

	"github.com/sourcegraph/jsonrpc2"
)

type handler struct{}

func (h *handler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	fmt.Println("OK SO FAR EVERYTHING SEEMS GOOD")
	switch req.Method {
	case "getStatus":
		h.getStatus(ctx, conn, req)
	default:
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{Code: 1, Message: fmt.Sprintf("Method %q not found", req.Method), Data: nil})
	}
}

func (h *handler) getStatus(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (string, error) {
	return "OK", nil
}
