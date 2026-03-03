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

func TestTsEnumValidationGeneratedAndWired(t *testing.T) {
	tmpDir := t.TempDir()
	g := NewTsGenerator(gen.Config{TsDir: tmpDir})

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
					{Name: "operators", Type: Type{Name: "sim_operator", Kind: KindEnum, IsList: true}},
				},
			},
		},
		Apis: []Api{
			{
				Name: "sim.set",
				Args: []ApiArg{
					{Name: "operator", Type: Type{Name: "sim_operator", Kind: KindEnum}},
					{Name: "operators", Type: Type{Name: "sim_operator", Kind: KindEnum, IsList: true}},
				},
				Result: Type{Name: "sim_operator", Kind: KindEnum},
			},
		},
	}

	if err := g.Generate(schema); err != nil {
		t.Fatalf("generate ts failed: %v", err)
	}

	enumContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "enum.ts"))
	if err != nil {
		t.Fatalf("read enum.ts failed: %v", err)
	}
	enumText := string(enumContent)
	assertContains(t, enumText, "export const IsSimOperator = (v: SimOperator): boolean => {")
	assertContains(t, enumText, "export const IsSimOperatorList = (v: SimOperator[]): boolean => {")

	structContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_sim.ts"))
	if err != nil {
		t.Fatalf("read struct_sim.ts failed: %v", err)
	}
	structText := string(structContent)
	assertContains(t, structText, "if (!_.IsSimOperatorList(s.operators as any)) return [s, new Error(\"get Sim Operators: invalid enum value\")];")
	assertContains(t, structText, "if ((s.operator as any) !== 0 && !_.IsSimOperator(s.operator as any)) return [s, new Error(\"get Sim Operator: invalid enum value\")];")
	assertContains(t, structText, "if (!_.IsSimOperatorList(s.operators as any)) return new Error(\"set Sim Operators: invalid enum value\");")
	assertContains(t, structText, "if (!_.IsSimOperator(s.operator as any)) return new Error(\"set Sim Operator: invalid enum value\");")

	rpcContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "rpc.ts"))
	if err != nil {
		t.Fatalf("read rpc.ts failed: %v", err)
	}
	rpcText := string(rpcContent)
	assertContains(t, rpcText, "if ((operator as any) !== 0 && !_.IsSimOperator(operator as any))")
	assertContains(t, rpcText, "if (!_.IsSimOperatorList(operators as any))")
	assertContains(t, rpcText, "if (err === null && (result as any) !== 0 && !_.IsSimOperator(result as any))")
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
	assertContains(t, text, "func CallFooPing(c *Client, ctx context.Context)")
	assertContains(t, text, "func doClient(c *Client, ctx context.Context, path string, body []byte)")
	assertNotContains(t, text, "func (c *Client)")
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
				Name: "sim_info",
				Fields: []StructField{
					{Name: "id", Type: Type{Name: "u32", Kind: KindBase}},
				},
			},
			{
				Name: "sim",
				Fields: []StructField{
					{Name: "operator", Type: Type{Name: "sim_operator", Kind: KindEnum}},
					{Name: "operators", Type: Type{Name: "sim_operator", Kind: KindEnum, IsList: true}},
					{Name: "info", Type: Type{Name: "sim_info", Kind: KindStruct}},
					{Name: "infos", Type: Type{Name: "sim_info", Kind: KindStruct, IsList: true}},
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
	assertContains(t, enumText, "func IsSimOperator(v SimOperator) bool")
	assertContains(t, enumText, "func IsSimOperatorList(v SimOperatorList) bool")
	assertContains(t, enumText, "func GetSimOperator(buf *bytes.Buffer) (SimOperator, error)")
	assertContains(t, enumText, "func SetSimOperator(buf *bytes.Buffer, v SimOperator) error")
	assertNotContains(t, enumText, "func (v SimOperatorList)")

	structContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "struct_sim.go"))
	if err != nil {
		t.Fatalf("read struct_sim.go failed: %v", err)
	}
	structText := string(structContent)
	assertContains(t, structText, "func ValidateSim(s *Sim) error")
	assertContains(t, structText, "if s.Operator != 0 && !IsSimOperator(s.Operator)")
	assertContains(t, structText, "if err := ValidateSim(s); err != nil")
	assertContains(t, structText, "func GetSim(buf *bytes.Buffer, s *Sim) error")
	assertContains(t, structText, "func ReadSim(buf *bytes.Buffer) (*Sim, error)")
	assertContains(t, structText, "func SetSim(buf *bytes.Buffer, s *Sim) error")
	assertContains(t, structText, "func GetSimList(buf *bytes.Buffer) (SimList, error)")
	assertContains(t, structText, "func SetSimList(buf *bytes.Buffer, v SimList) error")
	assertContains(t, structText, "func EqSim(a, b *Sim) bool")
	assertContains(t, structText, "func EqSimList(a, b SimList) bool")
	assertContains(t, structText, "if !EqSimOperatorList(a.Operators, b.Operators) { return false }")
	assertContains(t, structText, "if !EqSimInfo(a.Info, b.Info) { return false }")
	assertContains(t, structText, "if !EqSimInfoList((SimInfoList)(a.Infos), (SimInfoList)(b.Infos)) { return false }")
	assertNotContains(t, structText, "func (s *Sim) Get(")
	assertNotContains(t, structText, "func (s *Sim) Set(")
	assertNotContains(t, structText, "func (s *Sim) Eq(")
	assertNotContains(t, structText, ".Eq(")
	assertNotContains(t, structText, ".Get(")
	assertNotContains(t, structText, ".Set(")

	apiContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "api._.go"))
	if err != nil {
		t.Fatalf("read api._.go failed: %v", err)
	}
	apiText := string(apiContent)
	assertContains(t, apiText, "if !IsSimOperator(operator)")
	assertContains(t, apiText, "if !IsSimOperatorList(operators)")
	assertContains(t, apiText, "if err := ValidateSim(sim); err != nil")
	assertContains(t, apiText, "func parseRequest(w http.ResponseWriter, r *http.Request, args ...Getter) bool")
	assertContains(t, apiText, "func sendResponse(w http.ResponseWriter, setter Setter)")

	typeContent, err := os.ReadFile(filepath.Join(tmpDir, "sb", "type.go"))
	if err != nil {
		t.Fatalf("read type.go failed: %v", err)
	}
	typeText := string(typeContent)
	assertContains(t, typeText, "type Setter func(*bytes.Buffer) error")
	assertContains(t, typeText, "type Getter func(*bytes.Buffer) error")
	assertNotContains(t, typeText, "type Serializable interface")
	assertNotContains(t, typeText, "type Deserializable interface")
	assertNotContains(t, typeText, "func (v ")
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
