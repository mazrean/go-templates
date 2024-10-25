package router

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"connectrpc.com/connect"
	"connectrpc.com/grpchealth"
	"connectrpc.com/grpcreflect"
	"connectrpc.com/otelconnect"
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

	otelInterceptor, err := otelconnect.NewInterceptor()
	if err != nil {
		slog.Error(
			"failed to create otel interceptor",
			slog.String("error", err.Error()),
		)
	}

	logInterceptor := NewLogInterceptor()

	path, handler := protogenconnect.NewExampleServiceHandler(
		r.example,
		connect.WithInterceptors(otelInterceptor, logInterceptor),
	)
	mux.Handle(path, handler)

	serviceNames := []string{
		protogenconnect.ExampleServiceName,
	}
	mux.Handle(grpchealth.NewHandler(grpchealth.NewStaticChecker(serviceNames...)))
	mux.Handle(grpcreflect.NewHandlerV1(grpcreflect.NewStaticReflector(serviceNames...)))

	return http.ListenAndServe(
		addr,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

type logInterceptor struct{}

func NewLogInterceptor() connect.Interceptor {
	return &logInterceptor{}
}

func (l *logInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		headerMap := req.Header()
		headers := make([]any, 0, len(headerMap))
		for key, header := range headerMap {
			headers = append(headers, slog.String(key, strings.Join(header, ",")))
		}

		queryMap := req.Peer().Query
		queries := make([]any, 0, len(queryMap))
		for key, query := range queryMap {
			queries = append(queries, slog.String(key, strings.Join(query, ",")))
		}

		res, err := next(ctx, req)

		args := []any{
			slog.String("procedure", req.Spec().Procedure),
			slog.String("protocol", req.Peer().Protocol),
			slog.String("addr", req.Peer().Addr),
			slog.Group("request",
				slog.String("method", req.HTTPMethod()),
				slog.Group("header", headers...),
				slog.Group("query", queries...),
			),
		}
		if err != nil {
			args = append(args, slog.String("error", err.Error()))
			slog.ErrorContext(
				ctx, "unary",
				args...,
			)
		} else {
			slog.InfoContext(
				ctx, "unary",
				args...,
			)
		}

		return res, err
	}
}

func (l *logInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		conn := next(ctx, spec)

		requestHeaderMap := conn.RequestHeader()
		requestHeaders := make([]any, 0, len(requestHeaderMap))
		for key, header := range requestHeaderMap {
			requestHeaders = append(requestHeaders, slog.String(key, strings.Join(header, ",")))
		}

		requestQueryMap := conn.Peer().Query
		requestQueries := make([]any, 0, len(requestQueryMap))
		for key, query := range requestQueryMap {
			requestQueries = append(requestQueries, slog.String(key, strings.Join(query, ",")))
		}

		responseHeaderMap := conn.ResponseHeader()
		responseHeaders := make([]any, 0, len(responseHeaderMap))
		for key, header := range responseHeaderMap {
			responseHeaders = append(responseHeaders, slog.String(key, strings.Join(header, ",")))
		}

		responseTrailerMap := conn.ResponseTrailer()
		responseTrailers := make([]any, 0, len(responseTrailerMap))
		for key, trailer := range responseTrailerMap {
			responseTrailers = append(responseTrailers, slog.String(key, strings.Join(trailer, ",")))
		}

		slog.InfoContext(
			ctx, "streaming client",
			slog.String("procedure", spec.Procedure),
			slog.String("protocol", conn.Peer().Protocol),
			slog.String("addr", conn.Peer().Addr),
			slog.Group("request",
				slog.Group("header", requestHeaders...),
				slog.Group("query", requestQueries...),
			),
			slog.Group("response",
				slog.Group("header", responseHeaders...),
				slog.Group("trailer", responseTrailers...),
			),
		)

		return conn
	}
}

func (l *logInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		requestHeaderMap := conn.RequestHeader()
		requestHeaders := make([]any, 0, len(requestHeaderMap))
		for key, header := range requestHeaderMap {
			requestHeaders = append(requestHeaders, slog.String(key, strings.Join(header, ",")))
		}

		requestQueryMap := conn.Peer().Query
		requestQueries := make([]any, 0, len(requestQueryMap))
		for key, query := range requestQueryMap {
			requestQueries = append(requestQueries, slog.String(key, strings.Join(query, ",")))
		}

		slog.InfoContext(
			ctx, "streaming handler start",
			slog.String("procedure", conn.Spec().Procedure),
			slog.String("protocol", conn.Peer().Protocol),
			slog.String("addr", conn.Peer().Addr),
			slog.Group("request",
				slog.Group("header", requestHeaders...),
				slog.Group("query", requestQueries...),
			),
		)

		err := next(ctx, conn)
		if err != nil {
			slog.Error("streaming handler error", slog.String("error", err.Error()))
		}

		slog.InfoContext(ctx, "streaming handler end")

		return err
	}
}
