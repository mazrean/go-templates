//go:build tools

package main

import (
	_ "connectrpc.com/connect/cmd/protoc-gen-connect-go"
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
