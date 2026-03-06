package internal

import "strings"

func tplNoteLines(note string) []string {
	note = strings.ReplaceAll(note, "\r\n", "\n")
	note = strings.ReplaceAll(note, "\r", "\n")
	note = strings.TrimRight(note, "\n")
	if note == "" {
		return nil
	}
	return strings.Split(note, "\n")
}

func tplRenderLineComment(prefix, note string) string {
	lines := tplNoteLines(note)
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

func tplRenderLineCommentWithHead(prefix, head, note string) string {
	lines := tplNoteLines(note)
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

func tplRenderMarkdownInline(note string) string {
	lines := tplNoteLines(note)
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "<br>")
}

func tplRenderMarkdownQuote(note string) string {
	lines := tplNoteLines(note)
	if len(lines) == 0 {
		return ""
	}

	var b strings.Builder
	for i, line := range lines {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString("> ")
		b.WriteString(line)
	}
	return b.String()
}
