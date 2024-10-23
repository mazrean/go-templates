package router

import (
	"net/http"

	"github.com/mazrean/go-templates/connectrpc/internal/router/protogen/protobuf/protogenconnect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install connectrpc.com/connect/cmd/protoc-gen-connect-go
//go:generate go run github.com/bufbuild/buf/cmd/buf generate

type Router struct {
	example *Example
}

func NewRouter(
	example *Example,
) *Router {
	return &Router{
		example: example,
	}
}

func (r *Router) Run(addr string) error {
	mux := http.NewServeMux()

	path, handler := protogenconnect.NewExampleServiceHandler(r.example)
	mux.Handle(path, handler)

	return http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
