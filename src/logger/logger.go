package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	traceIdKey = "traceID"
)

// Logger is a zap logger provider struct
type Logger struct {
	*zap.Logger
}

func NewLogger(c *LoggerConfig) (*Logger, error) {
	l, err := GetLogger(c)
	if err != nil {
		return nil, err
	}

	return &Logger{l}, nil
}

// InfoCtx logs new InfoLevel messages considering current ctx
func (log *Logger) InfoCtx(ctx context.Context, msg string, fields ...*Field) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, getZapFieldsWithCtx(ctx, fields...)...)
}

// ErrorCtx logs new ErrorLevel messages considering current ctx
func (log *Logger) ErrorCtx(ctx context.Context, msg string, fields ...*Field) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, getZapFieldsWithCtx(ctx, fields...)...)
}

// DebugCtx logs new DebugLevel messages considering current ctx
func (log *Logger) DebugCtx(ctx context.Context, msg string, fields ...*Field) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, getZapFieldsWithCtx(ctx, fields...)...)
}

// WarnCtx logs new WarnLevel messages considering current ctx
func (log *Logger) WarnCtx(ctx context.Context, msg string, fields ...*Field) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, getZapFieldsWithCtx(ctx, fields...)...)
}

// Info logs new InfoLevel messages without considering current ctx
func (log *Logger) Info(msg string, fields ...*Field) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Info(msg, getZapFields(fields)...)
}

// Error logs new ErrorLevel messages without considering current ctx
func (log *Logger) Error(msg string, fields ...*Field) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Error(msg, getZapFields(fields)...)
}

// Debug logs new DebugLevel messages without considering current ctx
func (log *Logger) Debug(msg string, fields ...*Field) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Debug(msg, getZapFields(fields)...)
}

// Warn logs new WarnLevel messages without considering current ctx
func (log *Logger) Warn(msg string, fields ...*Field) {
	log.Logger.WithOptions(zap.AddCallerSkip(1)).Warn(msg, getZapFields(fields)...)
}

// Field is a context passed in Logger 
type Field struct {
	Key   string
	Value any
}

// NewFiled creates a new instance of any Field
func NewField(key string, value any) *Field {
	return &Field{
		Key:   key,
		Value: value,
	}
}

// NewErrorFiled creates a new instance of error Field
func NewErrorField(err error) *Field {
	return &Field{
		Key:   "error",
		Value: err,
	}
}

// getZapFieldsWithCtx gets zap.Fields array with tracingIDs
func getZapFieldsWithCtx(ctx context.Context, fields ...*Field) []zap.Field {
	zapFields := getZapFields(fields)
	zapFields = append(zapFields, zap.String(traceIdKey, getTracingIDs(ctx)))
	return zapFields
}

// getZapFields gets zap.Fields array without tracingIDs
func getZapFields(fields []*Field) []zap.Field {
	var zapFields []zap.Field
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}

	return zapFields
}

// getTracingIDs returns string tracingID from given ctx
func getTracingIDs(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}
