# syntax=docker/dockerfile:1
FROM golang:1.23.2 AS build

WORKDIR /app
RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go install github.com/grpc-ecosystem/grpc-health-probe
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o /bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=build /bin/app /bin/app
COPY --from=build /go/bin/grpc-health-probe /bin/grpc-health-probe

ENTRYPOINT ["/bin/app"]
