package log

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
)

// LogsProducer creates logs with tracing IDs and works with zap's interceptor
func LogsProducer(ctx context.Context, msg string, level zapcore.Level, code codes.Code, err error, duration zapcore.Field) {
	ctxzap.Extract(ctx).Check(level, msg).Write(
		zap.Error(err),
		zap.String("grpc.code", code.String()),
		duration,
		zap.String("traceID", getTracingId(ctx)),
	)
}

func getTracingId(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
} 