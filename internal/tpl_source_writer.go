package internal

import (
	"fmt"
	"strings"
)

type sourceWriter struct {
	b strings.Builder
}

func (w *sourceWriter) Line(line string) {
	w.b.WriteString(line)
	w.b.WriteByte('\n')
}

func (w *sourceWriter) Linef(format string, args ...any) {
	fmt.Fprintf(&w.b, format, args...)
	w.b.WriteByte('\n')
}

func (w *sourceWriter) Blank() {
	w.b.WriteByte('\n')
}

func (w *sourceWriter) Write(text string) {
	if text == "" {
		return
	}
	w.b.WriteString(text)
}

func (w *sourceWriter) WriteLineComment(prefix, note string) {
	text := tplRenderLineComment(prefix, note)
	if text == "" {
		return
	}
	w.Line(text)
}

func (w *sourceWriter) WriteLineCommentWithHead(prefix, head, note string) {
	text := tplRenderLineCommentWithHead(prefix, head, note)
	if text == "" {
		return
	}
	w.Line(text)
}

func (w *sourceWriter) String() string {
	return w.b.String()
}

func joinWithComma(items []string) string {
	return strings.Join(items, ", ")
}
