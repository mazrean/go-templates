// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: protobuf/example.proto

package protogenconnect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	protobuf "github.com/mazrean/go-templates/connectrpc/internal/router/protogen/protobuf"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// ExampleServiceName is the fully-qualified name of the ExampleService service.
	ExampleServiceName = "example.ExampleService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// ExampleServiceGreetProcedure is the fully-qualified name of the ExampleService's Greet RPC.
	ExampleServiceGreetProcedure = "/example.ExampleService/Greet"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	exampleServiceServiceDescriptor     = protobuf.File_protobuf_example_proto.Services().ByName("ExampleService")
	exampleServiceGreetMethodDescriptor = exampleServiceServiceDescriptor.Methods().ByName("Greet")
)

// ExampleServiceClient is a client for the example.ExampleService service.
type ExampleServiceClient interface {
	Greet(context.Context, *connect.Request[protobuf.GreetRequest]) (*connect.Response[protobuf.GreetResponse], error)
}

// NewExampleServiceClient constructs a client for the example.ExampleService service. By default,
// it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and
// sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC()
// or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewExampleServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) ExampleServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &exampleServiceClient{
		greet: connect.NewClient[protobuf.GreetRequest, protobuf.GreetResponse](
			httpClient,
			baseURL+ExampleServiceGreetProcedure,
			connect.WithSchema(exampleServiceGreetMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// exampleServiceClient implements ExampleServiceClient.
type exampleServiceClient struct {
	greet *connect.Client[protobuf.GreetRequest, protobuf.GreetResponse]
}

// Greet calls example.ExampleService.Greet.
func (c *exampleServiceClient) Greet(ctx context.Context, req *connect.Request[protobuf.GreetRequest]) (*connect.Response[protobuf.GreetResponse], error) {
	return c.greet.CallUnary(ctx, req)
}

// ExampleServiceHandler is an implementation of the example.ExampleService service.
type ExampleServiceHandler interface {
	Greet(context.Context, *connect.Request[protobuf.GreetRequest]) (*connect.Response[protobuf.GreetResponse], error)
}

// NewExampleServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewExampleServiceHandler(svc ExampleServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	exampleServiceGreetHandler := connect.NewUnaryHandler(
		ExampleServiceGreetProcedure,
		svc.Greet,
		connect.WithSchema(exampleServiceGreetMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/example.ExampleService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case ExampleServiceGreetProcedure:
			exampleServiceGreetHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedExampleServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedExampleServiceHandler struct{}

func (UnimplementedExampleServiceHandler) Greet(context.Context, *connect.Request[protobuf.GreetRequest]) (*connect.Response[protobuf.GreetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("example.ExampleService.Greet is not implemented"))
}
