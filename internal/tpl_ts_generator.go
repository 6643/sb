package internal

import (
	"path/filepath"
	"strings"
)

type TsGenerator struct {
	Config Config
}

func NewTsGenerator(cfg Config) *TsGenerator {
	return &TsGenerator{Config: cfg}
}

func (g *TsGenerator) getTsType(t TplType) string {
	switch t.Name {
	case "i8", "u8", "i16", "u16", "i32", "u32", "f32", "f64":
		return "number"
	case "i64", "u64":
		return "bigint"
	case "bool":
		return "boolean"
	case "text":
		return "string"
	case "bin":
		return "Uint8Array"
	}
	return PascalCase(t.Name)
}

func (g *TsGenerator) getTsLogicType(t TplType) string {
	if t.Name == "nil" {
		return "void"
	}
	base := g.getTsType(TplType{Name: t.Name, Kind: t.Kind})
	if !t.IsList {
		return base
	}
	return base + "[]"
}

func (g *TsGenerator) getTsValue(name string) string {
	switch name {
	case "text":
		return `""`
	case "bin":
		return "new Uint8Array(0)"
	case "bool":
		return "false"
	case "i64", "u64":
		return "0n"
	default:
		return "0"
	}
}

func (g *TsGenerator) Generate(schema *TplSchema) error {
	targetDir := filepath.Join(g.Config.TsDir, "sb")
	dir, err := newGeneratedDir(targetDir)
	if err != nil {
		return err
	}
	if err := dir.WriteAll(
		generatedFile{RelativePath: "type.ts", Data: []byte(renderTsRuntimeSource()), Perm: 0644},
		generatedFile{RelativePath: "enum.ts", Data: []byte(g.renderTsEnumFile(schema.Enums)), Perm: 0644},
	); err != nil {
		return err
	}
	structFiles := make([]string, 0, len(schema.Structs))
	for _, s := range schema.Structs {
		filename := "struct_" + SnakeCase(s.Name) + ".ts"
		structFiles = append(structFiles, strings.TrimSuffix(filename, ".ts"))
		if err := dir.Write(filename, []byte(g.renderTsStructFile(s)), 0644); err != nil {
			return err
		}
	}
	allFiles := append([]string{"enum"}, structFiles...)
	if err := dir.Write("_.ts", []byte(g.renderTsIndexFile(allFiles)), 0644); err != nil {
		return err
	}
	if len(schema.Apis) == 0 {
		return nil
	}
	if err := dir.WriteAll(
		generatedFile{RelativePath: "rpc.ts", Data: []byte(g.renderTsRPCFile(schema.Apis)), Perm: 0644},
		generatedFile{RelativePath: "rpc_smoke.test.ts", Data: []byte(g.renderTsSmokeTest(schema.Apis)), Perm: 0644},
	); err != nil {
		return err
	}
	return nil
}
