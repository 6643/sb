package tplgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sb/internal/gen"
)

func TestTsRPCUsesSetterClosuresForAllArgKinds(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewTsGenerator(gen.Config{TsDir: tmpDir})

	schema := &Schema{
		Structs: []Struct{
			{Name: "sim_info"},
		},
		Apis: []Api{
			{
				Name: "foo.test",
				Args: []ApiArg{
					{Name: "nums", Type: Type{Name: "u16", Kind: KindBase, IsList: true}},
					{Name: "id", Type: Type{Name: "u64", Kind: KindBase}},
					{Name: "score", Type: Type{Name: "i64", Kind: KindBase}},
					{Name: "infos", Type: Type{Name: "sim_info", Kind: KindStruct, IsList: true}},
					{Name: "ok", Type: Type{Name: "bool", Kind: KindBase}},
				},
				Result: Type{Name: "nil", Kind: KindNil},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.ts"))
	if err != nil {
		t.Fatalf("read rpc.ts failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "_.setU16List(buf, nums)")
	assertContains(t, text, "_.setU64(buf, id)")
	assertContains(t, text, "_.setI64(buf, score)")
	assertContains(t, text, "_.setSimInfoList(buf, infos as any)")
	assertContains(t, text, "_.setBool(buf, ok)")
	assertContains(t, text, "import * as _ from \"./_\"")
	assertNotContains(t, text, "_.setU16List(nums)")
	assertNotContains(t, text, "from \"./_.ts\"")

	indexContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "_.ts"))
	if err != nil {
		t.Fatalf("read _.ts failed: %v", err)
	}
	indexText := string(indexContent)
	assertContains(t, indexText, "export * from \"./type\"")
	assertNotContains(t, indexText, "export * from \"./type.ts\"")
}

func TestGoRPCHandlesReadAllErrorAndUsesClientTimeout(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(gen.Config{GoDir: tmpDir})

	schema := &Schema{
		Apis: []Api{
			{
				Name:   "foo.ping",
				Result: Type{Name: "nil", Kind: KindNil},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.go"))
	if err != nil {
		t.Fatalf("read rpc.go failed: %v", err)
	}
	text := string(content)

	assertContains(t, text, "context.WithTimeout(ctx, c.Timeout)")
	assertContains(t, text, "b, readErr := io.ReadAll(resp.Body)")
	assertContains(t, text, "if readErr != nil")
}

func TestGoApiHandlersValidateEnumAndStructArgs(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewGoGenerator(gen.Config{GoDir: tmpDir})

	schema := &Schema{
		Enums: []Enum{
			{
				Name: "sim_operator",
				Children: []EnumChild{
					{ID: 2, Name: "zz"},
					{ID: 3, Name: "lt"},
				},
			},
		},
		Structs: []Struct{
			{
				Name: "sim",
				Fields: []StructField{
					{Name: "operator", Type: Type{Name: "sim_operator", Kind: KindEnum}},
				},
			},
		},
		Apis: []Api{
			{
				Name: "sim.set",
				Args: []ApiArg{
					{Name: "operator", Type: Type{Name: "sim_operator", Kind: KindEnum}},
					{Name: "operators", Type: Type{Name: "sim_operator", Kind: KindEnum, IsList: true}},
					{Name: "sim", Type: Type{Name: "sim", Kind: KindStruct}},
				},
				Result: Type{Name: "nil", Kind: KindNil},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate go failed: %v", err)
	}

	enumContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "enum.go"))
	if err != nil {
		t.Fatalf("read enum.go failed: %v", err)
	}
	enumText := string(enumContent)
	assertContains(t, enumText, "func IsSimOperatorValid(v SimOperator) bool")
	assertContains(t, enumText, "func IsSimOperatorListValid(v SimOperatorList) bool")

	structContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_sim.go"))
	if err != nil {
		t.Fatalf("read struct_sim.go failed: %v", err)
	}
	structText := string(structContent)
	assertContains(t, structText, "func (s *Sim) Validate() error")
	assertContains(t, structText, "if s.Operator != 0 && !IsSimOperatorValid(s.Operator)")
	assertContains(t, structText, "if err := s.Validate(); err != nil")

	apiContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "api._.go"))
	if err != nil {
		t.Fatalf("read api._.go failed: %v", err)
	}
	apiText := string(apiContent)
	assertContains(t, apiText, "if !IsSimOperatorValid(SimOperator(operator))")
	assertContains(t, apiText, "if !IsSimOperatorListValid(operators)")
	assertContains(t, apiText, "if err := (&sim).Validate(); err != nil")
}

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
