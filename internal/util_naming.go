package internal

import (
	"strings"
	"unicode"
)

func SnakeCase(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}

	var b strings.Builder
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1])) {
			b.WriteByte('_')
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

func PascalCase(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == '_' || r == '-' })
	if len(parts) == 0 {
		return ""
	}

	var b strings.Builder
	for _, p := range parts {
		r := []rune(p)
		if len(r) == 0 {
			continue
		}
		b.WriteRune(unicode.ToUpper(r[0]))
		if len(r) > 1 {
			b.WriteString(string(r[1:]))
		}
	}
	return b.String()
}

func CamelCase(s string) string {
	p := PascalCase(s)
	r := []rune(p)
	if len(r) == 0 {
		return ""
	}
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func APIGroup(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		return parts[0]
	}
	return "api"
}
