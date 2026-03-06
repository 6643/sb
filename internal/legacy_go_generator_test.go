package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateLegacyGoRendersMultilineNotesSafely(t *testing.T) {
	tmpDir := t.TempDir()
	schema := &IRSchema{
		Enums: []IREnum{
			{
				Name: "order_status",
				Note: "订单状态\n待处理说明",
				Members: []IREnumMember{
					{ID: 0, Name: "pending", Note: "待处理\n可继续"},
				},
			},
		},
		Structs: []IRStruct{
			{
				Name: "sim_info",
				Note: "SIM信息\n用于展示",
				Fields: []IRField{
					{
						Name: "id",
						Type: IRType{Name: "u32", Kind: IRKindBase},
						Tag:  "_id",
						Note: "主键\n数据库ID",
					},
				},
			},
		},
	}

	if err := GenerateLegacyGo(schema, Config{GoDir: tmpDir, GoTag: "json"}); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	enumText := legacyGoReadText(t, filepath.Join(tmpDir, "sb", "enum.go"))
	legacyGoAssertContains(t, enumText, "// OrderStatus 订单状态\n// 待处理说明\ntype OrderStatus uint8")
	legacyGoAssertContains(t, enumText, "\t// 待处理\n\t// 可继续\n\tOrderStatusPending OrderStatus = 0")
	legacyGoAssertNotContains(t, enumText, "\n待处理说明\ntype OrderStatus")
	legacyGoAssertNotContains(t, enumText, "OrderStatusPending OrderStatus = 0 // 待处理")

	structText := legacyGoReadText(t, filepath.Join(tmpDir, "sb", "struct_sim_info.go"))
	legacyGoAssertContains(t, structText, "// SimInfo SIM信息\n// 用于展示\ntype SimInfo struct {")
	legacyGoAssertContains(t, structText, "\t// 主键\n\t// 数据库ID\n\tId uint32 `json:\"_id\"`")
	legacyGoAssertNotContains(t, structText, "\n数据库ID\n\tId uint32")
}

func legacyGoReadText(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s failed: %v", path, err)
	}
	return string(content)
}

func legacyGoAssertContains(t *testing.T, text, want string) {
	t.Helper()
	if strings.Contains(text, want) {
		return
	}
	t.Fatalf("missing text %q", want)
}

func legacyGoAssertNotContains(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		return
	}
	t.Fatalf("unexpected text %q", want)
}
