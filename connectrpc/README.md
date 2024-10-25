# Go Connnect Template

This is a template for a Go Connect RPC service.

## Getting Started
```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/mazrean/go-templates/connectrpc your.domain/myprog
```

## Start the development server
```bash
docker compose watch
```

## After Updating the Proto File
```bash
go generate ./...
```
