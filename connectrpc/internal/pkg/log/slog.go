package log

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

func init() {
	Setup(true)
}

func Setup(isDebug bool) {
	var level slog.Level
	if isDebug {
		level = slog.LevelDebug
	} else {
		level = slog.LevelInfo
	}

	logger := slog.New(NewHandler(level))
	slog.SetDefault(logger)
}

type Handler struct {
	slog.Handler
}

func NewHandler(level slog.Level) *Handler {
	return &Handler{
		Handler: slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}),
	}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	return h.Handler.Handle(ctx, r)
}
