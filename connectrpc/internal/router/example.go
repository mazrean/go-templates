package router

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	protogen "github.com/mazrean/go-templates/connectrpc/internal/router/protogen/protobuf"
	"github.com/mazrean/go-templates/connectrpc/internal/router/protogen/protobuf/protogenconnect"
)

var _ protogenconnect.ExampleServiceHandler = &Example{}

type Example struct{}

func NewExample() *Example {
	return &Example{}
}

func (g *Example) Greet(ctx context.Context, req *connect.Request[protogen.GreetRequest]) (*connect.Response[protogen.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&protogen.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})

	return res, nil
}
