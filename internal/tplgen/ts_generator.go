package tplgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sb/internal/gen"
	"sb/internal/util"
	"strings"
	"text/template"
)

type TsGenerator struct {
	Config   gen.Config
	FuncMap  template.FuncMap
	tplCache map[string]*template.Template
}

func NewTsGenerator(cfg gen.Config) *TsGenerator {
	g := &TsGenerator{Config: cfg, tplCache: make(map[string]*template.Template)}
	g.FuncMap = template.FuncMap{
		"PascalCase":  util.PascalCase,
		"SnakeCase":   util.SnakeCase,
		"CamelCase":   util.CamelCase,
		"TsType":      g.getTsType,
		"TsLogicType": g.getTsLogicType,
		"TsValue":     g.getTsValue,
		"IsBaseType":  func(t Type) bool { return t.Kind == KindBase },
		"IsEnum":      func(t Type) bool { return t.Kind == KindEnum },
		"IsStruct":    func(t Type) bool { return t.Kind == KindStruct },
		"IsList":      func(t Type) bool { return t.IsList },
	}
	return g
}

func (g *TsGenerator) getTsType(t Type) string {
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
	return util.PascalCase(t.Name)
}

func (g *TsGenerator) getTsLogicType(t Type) string {
	if t.Name == "nil" {
		return "void"
	}

	base := g.getTsType(Type{Name: t.Name, Kind: t.Kind})
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

func (g *TsGenerator) Generate(schema *Schema) error {
	targetDir := filepath.Join(g.Config.TsDir, "sb")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create ts target dir %s: %w", targetDir, err)
	}

	typeTS, err := tplFS.ReadFile("_tpl/type.ts")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(targetDir, "type.ts"), typeTS, 0644); err != nil {
		return err
	}

	if err := g.executeTemplate("_tpl/ts.enum.tpl", filepath.Join(targetDir, "enum.ts"), map[string]any{
		"Enums": schema.Enums,
	}); err != nil {
		return err
	}

	structFiles := make([]string, 0, len(schema.Structs))
	for _, s := range schema.Structs {
		filename := "struct_" + util.SnakeCase(s.Name) + ".ts"
		structFiles = append(structFiles, strings.TrimSuffix(filename, ".ts"))
		path := filepath.Join(targetDir, filename)
		if err := g.executeTemplate("_tpl/ts.struct.tpl", path, s); err != nil {
			return err
		}
	}

	allFiles := append([]string{"enum"}, structFiles...)
	if err := g.executeTemplate("_tpl/ts._.tpl", filepath.Join(targetDir, "_.ts"), allFiles); err != nil {
		return err
	}

	if len(schema.Apis) == 0 {
		return nil
	}

	if err := g.executeTemplate("_tpl/ts.rpc.tpl", filepath.Join(targetDir, "rpc.ts"), map[string]any{
		"Apis": schema.Apis,
	}); err != nil {
		return err
	}
	return nil
}

func (g *TsGenerator) executeTemplate(tplPath, destPath string, data any) error {
	tpl, ok := g.tplCache[tplPath]
	if !ok {
		tplContent, err := tplFS.ReadFile(tplPath)
		if err != nil {
			return err
		}

		tpl, err = template.New(filepath.Base(tplPath)).Funcs(g.FuncMap).Parse(string(tplContent))
		if err != nil {
			return err
		}

		g.tplCache[tplPath] = tpl
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return err
	}
	return os.WriteFile(destPath, buf.Bytes(), 0644)
}
