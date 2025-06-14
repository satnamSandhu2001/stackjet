package logger

import (
	"io"
)

// emits logs in real time to SSE w
func EmitLog(w io.Writer, message string) {
	io.WriteString(w, message+"\n")
	if flusher, ok := w.(interface{ Flush() }); ok {
		flusher.Flush()
	}
}
