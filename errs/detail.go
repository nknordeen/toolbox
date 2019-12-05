package errs

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Causer defines the interface for determining the error the caused an error.
type Causer interface {
	Cause() error
}

type detail struct {
	message string
	stack   []uintptr
	cause   error
}

// Error implements the error interface.
func (d detail) Error() string {
	return d.detail(true, true)
}

// Cause returns the cause of this error, if any.
func (d detail) Cause() error {
	return d.cause
}

// Format implements the fmt.Formatter interface.
//
// Supported formats:
//   - "%s"  Just the message
//   - "%q"  Just the message, but quoted
//   - "%v"  The message plus a stack trace, trimmed of golang runtime calls
//   - "%+v" The message plus a stack trace
func (d *detail) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		state.Write([]byte(d.detail(true, !state.Flag('+'))))
	case 's':
		state.Write([]byte(d.message))
	case 'q':
		fmt.Fprintf(state, "%q", d.message)
	}
}

// StackTrace returns the raw call stack pointers.
func (d *detail) StackTrace() []uintptr {
	return d.stack
}

func (d *detail) detail(includeMessage, trimRuntime bool) string {
	var buffer strings.Builder
	if includeMessage {
		buffer.WriteString(d.message)
	}
	frames := runtime.CallersFrames(d.stack)
	for {
		frame, more := frames.Next()
		if frame.Function != "" {
			if trimRuntime && (strings.HasPrefix(frame.Function, "runtime.") || strings.HasPrefix(frame.Function, "testing.")) {
				continue
			}
			buffer.WriteString("\n    [")
			buffer.WriteString(frame.Function)
			buffer.WriteString("] ")
			file := frame.File
			if i := strings.Index(file, "."); i != -1 {
				for i > 0 && file[i] != os.PathSeparator {
					i--
				}
				if i > 0 {
					file = file[i+1:]
				}
				if i = strings.LastIndexByte(file, os.PathSeparator); i != -1 {
					path := file[:i]
					offset := i + 1
					if i = strings.LastIndexByte(path, os.PathSeparator); i != -1 {
						if path[i+1:] == "_obj" {
							path = path[:i]
						}
					}
					if strings.HasPrefix(frame.Function, path) {
						file = file[offset:]
					}
				}
			}
			buffer.WriteString(file)
			buffer.WriteByte(':')
			buffer.WriteString(strconv.Itoa(frame.Line))
		}
		if !more {
			break
		}
	}
	if d.cause != nil {
		buffer.WriteString("\n  Caused by: ")
		if detailed, ok := d.cause.(*Error); ok {
			buffer.WriteString(detailed.Detail(trimRuntime))
		} else {
			buffer.WriteString(d.cause.Error())
		}
	}
	return buffer.String()
}
