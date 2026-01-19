package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"sb/internal/ast"
	"sb/internal/util"
	"text/template"
)

type TsGenerator struct {
	Config  Config
	FuncMap template.FuncMap
}

func NewTsGenerator(cfg Config) *TsGenerator {
	g := &TsGenerator{Config: cfg}
	g.FuncMap = template.FuncMap{
		"PascalCase":  util.PascalCase,
		"SnakeCase":   util.SnakeCase,
		"CamelCase":   util.CamelCase,
		"TsType":      g.getTsType,
		"TsValue":     g.getTsValue,
		"TsLogicType": g.getTsLogicType,
		"IsBaseType":  func(t ast.Type) bool { return t.Kind == ast.KindBase },
		"IsEnum":      func(t ast.Type) bool { return t.Kind == ast.KindEnum },
		"IsStruct":    func(t ast.Type) bool { return t.Kind == ast.KindStruct },
		"IsList":      func(t ast.Type) bool { return t.IsList },
	}
	return g
}

func (g *TsGenerator) getTsType(t ast.Type) string {
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

func (g *TsGenerator) getTsValue(name string) string {
	switch name {
	case "text": return `""`
	case "bin": return "new Uint8Array(0)"
	case "bool": return "false"
	case "i64", "u64": return "0n"
	default: return "0"
	}
}

func (g *TsGenerator) getTsLogicType(t ast.Type) string {
	suffix := ""
	if t.IsList { suffix = "[]" }
	return g.getTsType(t) + suffix
}

func (g *TsGenerator) Generate(schema *ast.Schema) error {
	targetDir := filepath.Join(g.Config.TsDir, "sb")
	os.MkdirAll(targetDir, 0755)

	// 0. Copy type.ts from embedded FS
	typeTs, err := g.Config.TplFS.ReadFile("_tpl/type.ts")
	if err != nil { return err }
	if err := os.WriteFile(filepath.Join(targetDir, "type.ts"), typeTs, 0644); err != nil { return err }

	// 1. Generate Enums
	if err := g.executeTemplate("_tpl/ts.enum.tpl", filepath.Join(targetDir, "enum.ts"), map[string]any{
		"Enums": schema.Enums,
	}); err != nil { return err }

	// 2. Generate Structs
	var structFiles []string
	for _, s := range schema.Structs {
		filename := "struct_" + util.SnakeCase(s.Name) + ".ts"
		structFiles = append(structFiles, filename)
		path := filepath.Join(targetDir, filename)
		if err := g.executeTemplate("_tpl/ts.struct.tpl", path, s); err != nil { return err }
	}

	// 3. Generate Index (_.ts)
	allFiles := append([]string{"enum.ts"}, structFiles...)
	if err := g.executeTemplate("_tpl/ts._.tpl", filepath.Join(targetDir, "_.ts"), allFiles); err != nil { return err }

	// 4. Generate RPC
	if len(schema.Apis) > 0 {
		if err := g.executeTemplate("_tpl/ts.rpc.tpl", filepath.Join(targetDir, "rpc.ts"), map[string]any{
			"Apis": schema.Apis,
		}); err != nil { return err }
	}

	return nil
}

func (g *TsGenerator) executeTemplate(tplPath, destPath string, data any) error {
	tplContent, err := g.Config.TplFS.ReadFile(tplPath)
	if err != nil { return err }
	tpl, err := template.New(filepath.Base(tplPath)).Funcs(g.FuncMap).Parse(string(tplContent))
	if err != nil { return err }
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil { return err }
	return os.WriteFile(destPath, buf.Bytes(), 0644)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil { return err }
	return os.WriteFile(dst, data, 0644)
}
