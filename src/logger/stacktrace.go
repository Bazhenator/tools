package logger

import (
	"fmt"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

// maxStackDepth is a max depth of printed stacktrace
const maxStackDepth = 99

// StackFrame represents complete info structure for logging 
type StackFrame struct {
	FuncName string
	Filename string
	Line     int
}

func NewStackFrame(framePC uintptr) *StackFrame {
	fn := runtime.FuncForPC(framePC)
	funcName := fn.Name()
	if fn == nil {
		return &StackFrame{
			FuncName: funcName,
			Filename: "unknown",
			Line:     -1,
		}
	}

	file, Line := fn.FileLine(framePC)
	return &StackFrame{
		FuncName: funcName,
		Filename: file,
		Line:     Line,
	}
}

// buildStacktraceLogString builds string value for given stacktrace
func buildStacktraceLogString(stack []uintptr) string {
	var builder strings.Builder
	builder.WriteString("Stack trace:\n")

	for i, frame := range stack {
		frameRepr := NewStackFrame(frame)
		builder.WriteString(fmt.Sprintf("%2d: %s at %s:%d", i, frameRepr.FuncName, frameRepr.Filename, frameRepr.Line))
		if i == len(stack)-1 {
			break
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// LogStacktrace logs stacktrace into zap from call moment to func main()
func LogStacktrace(l *zap.Logger) {
	// getting program counters' list of "interesting" stack frames
	pcs := make([]uintptr, maxStackDepth)
	n := runtime.Callers(2, pcs) // 2 skip runtime.Callers, LogStacktrace
	stack := pcs[:n-2]           // 2 skip runtime.goexit, runtime.main

	l.Error(buildStacktraceLogString(stack))
}