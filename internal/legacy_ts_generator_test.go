package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateLegacyTsRendersMultilineNotesSafely(t *testing.T) {
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
			},
		},
		APIs: []IRAPI{
			{
				Name:   "user.set_sim_info",
				Result: IRType{Name: "nil", Kind: IRKindNil},
				Note:   "设置sim信息\n无返回值",
			},
		},
	}

	if err := GenerateLegacyTs(schema, Config{TsDir: tmpDir}); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	enumText := legacyTsReadText(t, filepath.Join(tmpDir, "sb", "enum.ts"))
	legacyTsAssertContains(t, enumText, "// 订单状态\n// 待处理说明\nexport enum OrderStatus {")
	legacyTsAssertContains(t, enumText, "  // 待处理\n  // 可继续\n  Pending = 0,")
	legacyTsAssertNotContains(t, enumText, "\n待处理说明\nexport enum OrderStatus")
	legacyTsAssertNotContains(t, enumText, "Pending = 0, // 待处理")

	structText := legacyTsReadText(t, filepath.Join(tmpDir, "sb", "struct.ts"))
	legacyTsAssertContains(t, structText, "// SIM信息\n// 用于展示\nexport interface SimInfo {")

	rpcText := legacyTsReadText(t, filepath.Join(tmpDir, "sb", "rpc.ts"))
	legacyTsAssertContains(t, rpcText, "  // 设置sim信息\n  // 无返回值\n  public userSetSimInfo = async")
	legacyTsAssertNotContains(t, rpcText, "/** 设置sim信息")
}

func legacyTsReadText(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s failed: %v", path, err)
	}
	return string(content)
}

func legacyTsAssertContains(t *testing.T, text, want string) {
	t.Helper()
	if strings.Contains(text, want) {
		return
	}
	t.Fatalf("missing text %q", want)
}

func legacyTsAssertNotContains(t *testing.T, text, want string) {
	t.Helper()
	if !strings.Contains(text, want) {
		return
	}
	t.Fatalf("unexpected text %q", want)
}
