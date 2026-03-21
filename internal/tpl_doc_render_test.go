package internal

import (
	"strings"
	"testing"
)

func TestRenderDocUsesSbPublicImports(t *testing.T) {
	text := renderDoc(&TplSchema{
		Apis: []TplApi{{
			Name:   "user.get_count",
			Result: TplType{Name: "u8", Kind: TplKindBase},
		}},
	})

	if got := strings.Count(text, "\"your_project/sb\""); got != 2 {
		t.Fatalf("expected Go public import path twice, got %d", got)
	}
	assertContains(t, text, "import * as sb from \"./sb\";")
	assertContains(t, text, "client := sb.NewClient")
	assertNotContains(t, text, "\"your_project/go/sb\"")
}
