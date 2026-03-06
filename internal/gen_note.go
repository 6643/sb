package internal

import "strings"

func noteLines(note string) []string {
	note = strings.ReplaceAll(note, "\r\n", "\n")
	note = strings.ReplaceAll(note, "\r", "\n")
	note = strings.TrimRight(note, "\n")
	if note == "" {
		return nil
	}
	return strings.Split(note, "\n")
}

// RenderLineComment 将多行说明渲染为安全的逐行注释。
func RenderLineComment(prefix, note string) string {
	lines := noteLines(note)
	if len(lines) == 0 {
		return ""
	}

	var b strings.Builder
	for i, line := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(prefix)
		b.WriteString(line)
	}
	return b.String()
}

// RenderLineCommentWithHead 将首行附加标题, 后续行保持同前缀输出。
func RenderLineCommentWithHead(prefix, head, note string) string {
	lines := noteLines(note)
	if len(lines) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(prefix)
	b.WriteString(head)
	b.WriteString(lines[0])
	for _, line := range lines[1:] {
		b.WriteByte('\n')
		b.WriteString(prefix)
		b.WriteString(line)
	}
	return b.String()
}
