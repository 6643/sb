package internal

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestTsGeneratorWritesSchemaFile(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewTsGenerator(Config{TsDir: tmpDir})
	sbDir := filepath.Join(tmpDir, "sb")
	staleRuntimeFile := filepath.Join(sbDir, "runtime_text_io.ts")
	staleRuntimeTestFile := filepath.Join(sbDir, "runtime_text_io.test.ts")
	if err := os.MkdirAll(sbDir, 0o755); err != nil {
		t.Fatalf("create sb dir failed: %v", err)
	}
	if err := os.WriteFile(staleRuntimeFile, []byte("stale"), 0o644); err != nil {
		t.Fatalf("write stale runtime file failed: %v", err)
	}
	if err := os.WriteFile(staleRuntimeTestFile, []byte("keep"), 0o644); err != nil {
		t.Fatalf("write stale runtime test file failed: %v", err)
	}

	schema := &TplSchema{
		Enums: []TplEnum{
			{
				Name: "sim_operator",
				Children: []TplEnumChild{
					{ID: 2, Name: "zz"},
					{ID: 3, Name: "lt"},
				},
			},
		},
		Structs: []TplStruct{
			{
				Name: "sim_info",
				Fields: []TplStructField{
					{Name: "title", Type: TplType{Name: "text", Kind: TplKindBase}},
				},
			},
			{
				Name: "sim",
				Fields: []TplStructField{
					{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}},
					{Name: "name", Type: TplType{Name: "text", Kind: TplKindBase}},
					{Name: "operator", Type: TplType{Name: "sim_operator", Kind: TplKindEnum}},
					{Name: "ok", Type: TplType{Name: "bool", Kind: TplKindBase}},
					{Name: "tags", Type: TplType{Name: "u32", Kind: TplKindBase, IsList: true}},
					{Name: "titles", Type: TplType{Name: "text", Kind: TplKindBase, IsList: true}},
					{Name: "infos", Type: TplType{Name: "sim_info", Kind: TplKindStruct, IsList: true}},
				},
			},
		},
		Apis: []TplApi{
			{
				Name: "user.get_count",
				Args: []TplApiArg{
					{Name: "page", Type: TplType{Name: "u8", Kind: TplKindBase}},
				},
				Result: TplType{Name: "u8", Kind: TplKindBase},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	assertFileMissing(t, staleRuntimeFile)
	assertFileExists(t, staleRuntimeTestFile)

	baseText := mustReadFile(t, filepath.Join(sbDir, "runtime_base.ts"))
	assertContains(t, baseText, "export interface Buffer")
	assertContains(t, baseText, "export const Buffer = _BufferCtor as unknown as BufferCtor;")
	assertNotContains(t, baseText, "export * from \"./runtime_buffer\"")
	assertNotContains(t, baseText, "export * from \"./runtime_scalar\"")

	headerText := mustReadFile(t, filepath.Join(sbDir, "runtime_header.ts"))
	assertContains(t, headerText, "export const writeHeader")
	assertContains(t, headerText, "export const encodeBitmap")
	assertNotContains(t, headerText, "export * from \"./runtime_bitmap\"")
	assertNotContains(t, headerText, "export * from \"./runtime_state_header\"")

	textRuntimeText := mustReadFile(t, filepath.Join(sbDir, "runtime_text.ts"))
	assertContains(t, textRuntimeText, "export const textState")
	assertContains(t, textRuntimeText, "export const getText")
	assertNotContains(t, textRuntimeText, "export * from \"./runtime_text_utf8\"")
	assertNotContains(t, textRuntimeText, "export * from \"./runtime_text_io\"")

	binRuntimeText := mustReadFile(t, filepath.Join(sbDir, "runtime_bin.ts"))
	assertContains(t, binRuntimeText, "export const getBinLength")
	assertContains(t, binRuntimeText, "export const setBin")
	assertNotContains(t, binRuntimeText, "export * from \"./runtime_bin_state\"")
	assertNotContains(t, binRuntimeText, "export * from \"./runtime_bin_io\"")

	listRuntimeText := mustReadFile(t, filepath.Join(sbDir, "runtime_list.ts"))
	assertContains(t, listRuntimeText, "export const getDefaultList")
	assertContains(t, listRuntimeText, "export const getTextList")
	assertNotContains(t, listRuntimeText, "export * from \"./runtime_bitmap_list\"")
	assertNotContains(t, listRuntimeText, "export * from \"./runtime_item_list\"")

	enumRuntimeText := mustReadFile(t, filepath.Join(sbDir, "runtime_enum.ts"))
	assertContains(t, enumRuntimeText, "export interface EnumMeta<T extends number>")
	assertContains(t, enumRuntimeText, "export const getEnumList")
	assertNotContains(t, enumRuntimeText, "export * from \"./runtime_enum\"")

	structRuntimeText := mustReadFile(t, filepath.Join(sbDir, "runtime_struct.ts"))
	assertContains(t, structRuntimeText, "export interface StructMeta<T>")
	assertContains(t, structRuntimeText, "export const setStruct")
	assertContains(t, structRuntimeText, "export const boolField")
	assertContains(t, structRuntimeText, "export const defaultListField")
	assertNotContains(t, structRuntimeText, "export * from \"./runtime_struct_core\"")
	assertNotContains(t, structRuntimeText, "export * from \"./runtime_struct_field\"")

	assertRuntimeContractFilesExactly(t, sbDir, []string{
		"type.ts",
		"runtime_base.ts",
		"runtime_header.ts",
		"runtime_text.ts",
		"runtime_bin.ts",
		"runtime_list.ts",
		"runtime_struct.ts",
		"runtime_enum.ts",
		"runtime_core.ts",
		"runtime_meta.ts",
	})
	assertFileExists(t, filepath.Join(sbDir, "type.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "runtime_text_io.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "runtime_text_state.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "runtime_bin_length.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "runtime_bin_io.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "runtime_struct_scalar_field.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "runtime_struct_default_list_field.ts"))

	typeText := mustReadFile(t, filepath.Join(sbDir, "type.ts"))
	assertReExportOnly(t, typeText)
	assertContains(t, typeText, "export * from \"./runtime_core\"")
	assertContains(t, typeText, "export * from \"./runtime_meta\"")

	coreText := mustReadFile(t, filepath.Join(sbDir, "runtime_core.ts"))
	assertReExportOnly(t, coreText)
	assertContains(t, coreText, "export * from \"./runtime_base\"")
	assertContains(t, coreText, "export * from \"./runtime_header\"")
	assertContains(t, coreText, "export * from \"./runtime_text\"")
	assertContains(t, coreText, "export * from \"./runtime_bin\"")
	assertContains(t, coreText, "export * from \"./runtime_list\"")
	assertNoFineGrainedRuntimeRefs(t, typeText, []string{
		"./runtime_buffer",
		"./runtime_header_codec",
		"./runtime_text_io",
		"./runtime_bin_length",
		"./runtime_struct_scalar_field",
	})
	assertNoFineGrainedRuntimeRefs(t, coreText, []string{
		"./runtime_buffer",
		"./runtime_header_codec",
		"./runtime_text_io",
		"./runtime_bin_length",
		"./runtime_struct_scalar_field",
	})

	metaText := mustReadFile(t, filepath.Join(sbDir, "runtime_meta.ts"))
	assertReExportOnly(t, metaText)
	assertContains(t, metaText, "export * from \"./runtime_enum\"")
	assertContains(t, metaText, "export * from \"./runtime_struct\"")
	assertNotContains(t, metaText, "export interface EnumMeta<T extends number>")
	assertNotContains(t, metaText, "export interface StructMeta<T>")
	assertNotContains(t, metaText, "export const setStruct")
	assertNoFineGrainedRuntimeRefs(t, metaText, []string{
		"./runtime_struct_meta",
		"./runtime_struct_codec",
		"./runtime_struct_scalar_field",
		"./runtime_text_io",
		"./runtime_bin_length",
	})

	enumText := mustReadFile(t, filepath.Join(sbDir, "enum.ts"))
	assertNotContains(t, enumText, "from \"./type\"")
	assertContains(t, enumText, "from \"./runtime_core\"")
	assertContains(t, enumText, "from \"./runtime_meta\"")

	structText := mustReadFile(t, filepath.Join(sbDir, "struct_sim.ts"))
	assertNotContains(t, structText, "from \"./type\"")
	assertContains(t, structText, "from \"./runtime_core\"")
	assertContains(t, structText, "from \"./runtime_meta\"")

	rpcText := mustReadFile(t, filepath.Join(sbDir, "rpc.ts"))
	assertNotContains(t, rpcText, "from \"./type\"")
	assertContains(t, rpcText, "from \"./runtime_core\"")
	assertNotContains(t, rpcText, "from \"./runtime_meta\"")
}

func TestTsGeneratorRemovesStaleManagedSchemaFiles(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewTsGenerator(Config{TsDir: tmpDir})
	sbDir := filepath.Join(tmpDir, "sb")

	withDefs := &TplSchema{
		Enums: []TplEnum{
			{Name: "status", Children: []TplEnumChild{{ID: 0, Name: "ok"}}},
		},
		Structs: []TplStruct{
			{Name: "user", Fields: []TplStructField{{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}}}},
		},
		Apis: []TplApi{
			{Name: "get_count", Result: TplType{Name: "u32", Kind: TplKindBase}},
		},
	}
	changedDefs := &TplSchema{
		Structs: []TplStruct{
			{Name: "order", Fields: []TplStructField{{Name: "id", Type: TplType{Name: "u32", Kind: TplKindBase}}}},
		},
	}

	if err := g.Generate(withDefs); err != nil {
		t.Fatalf("generate schema with defs failed: %v", err)
	}
	manualTest := filepath.Join(sbDir, "rpc_manual.test.ts")
	if err := os.WriteFile(manualTest, []byte("keep"), 0o644); err != nil {
		t.Fatalf("write manual test failed: %v", err)
	}

	if err := g.Generate(changedDefs); err != nil {
		t.Fatalf("regenerate changed schema failed: %v", err)
	}

	assertFileMissing(t, filepath.Join(sbDir, "struct_user.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "rpc.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "rpc_smoke.test.ts"))
	assertFileMissing(t, filepath.Join(sbDir, "enum_smoke.test.ts"))
	assertFileExists(t, filepath.Join(sbDir, "struct_order.ts"))
	assertFileExists(t, manualTest)
	indexText := mustReadFile(t, filepath.Join(sbDir, "_.ts"))
	assertContains(t, indexText, "export * from \"./struct_order\"")
	assertNotContains(t, indexText, "export * from \"./struct_user\"")
}

func assertRuntimeContractFilesExactly(t *testing.T, dir string, want []string) {
	t.Helper()
	runtimeEntries, err := filepath.Glob(filepath.Join(dir, "runtime_*.ts"))
	if err != nil {
		t.Fatalf("glob runtime files failed: %v", err)
	}
	entries := append([]string{filepath.Join(dir, "type.ts")}, runtimeEntries...)
	got := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := filepath.Base(entry)
		if strings.HasSuffix(name, ".test.ts") {
			continue
		}
		got = append(got, name)
	}
	sort.Strings(got)
	want = append([]string(nil), want...)
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("unexpected runtime file count: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected runtime files: got=%v want=%v", got, want)
		}
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file exists: %s err=%v", path, err)
	}
}

func assertFileMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected file missing: %s err=%v", path, err)
	}
}

func mustReadFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s failed: %v", filepath.Base(path), err)
	}
	return string(content)
}

func assertReExportOnly(t *testing.T, text string) {
	t.Helper()
	for _, bad := range []string{
		"export interface ",
		"export const ",
		"export function ",
		"export class ",
		"export enum ",
		"export type ",
	} {
		assertNotContains(t, text, bad)
	}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "*/") {
			continue
		}
		if !strings.HasPrefix(line, "export * from ") {
			t.Fatalf("expected re-export-only file, got line: %q", line)
		}
	}
}

func assertNoFineGrainedRuntimeRefs(t *testing.T, text string, disallowed []string) {
	t.Helper()
	for _, ref := range disallowed {
		assertNotContains(t, text, ref)
	}
}
