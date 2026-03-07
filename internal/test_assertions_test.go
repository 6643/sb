package internal

import (
	"strings"
	"testing"
)

func assertContains(t *testing.T, text, want string) {
	t.Helper()
	if strings.Contains(text, want) {
		return
	}
	t.Fatalf("missing text %q", want)
}

func assertNotContains(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		return
	}
	t.Fatalf("unexpected text %q", want)
}
