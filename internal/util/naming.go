package util

import (
	"strings"
	"unicode"
)

func SnakeCase(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	var builder strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 && !unicode.IsUpper(rune(s[i-1])) && s[i-1] != '_' {
			builder.WriteRune('_')
		}
		builder.WriteRune(unicode.ToLower(r))
	}
	return builder.String()
}

func PascalCase(s string) string {
	s = strings.ReplaceAll(s, ".", "_")
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '.'
	})
	var builder strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			builder.WriteString(strings.ToUpper(part[:1]))
			builder.WriteString(part[1:])
		}
	}
	return builder.String()
}

func CamelCase(s string) string {
	pascal := PascalCase(s)
	if len(pascal) > 0 {
		return strings.ToLower(pascal[:1]) + pascal[1:]
	}
	return ""
}
