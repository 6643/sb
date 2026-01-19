package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sb/internal/ast"
	"sb/internal/util"
	"strings"
	"text/template"
	"math"
)

type GoGenerator struct {
	Config  Config
	FuncMap template.FuncMap
}

func NewGoGenerator(cfg Config) *GoGenerator {
	g := &GoGenerator{Config: cfg}
	g.FuncMap = template.FuncMap{
		"PascalCase":  util.PascalCase,
		"SnakeCase":   util.SnakeCase,
		"CamelCase":   util.CamelCase,
		"GoType":      g.getGoType,
		"GoValue":     g.getGoValue,
		"GoTag":       g.getGoTag,
		"GoLogicType": g.getGoLogicType,
		"GoRpcType":   g.getGoRpcType,
		"IsBaseType":  func(t ast.Type) bool { return t.Kind == ast.KindBase },
		"IsEnum":      func(t ast.Type) bool { return t.Kind == ast.KindEnum },
		"IsStruct":    func(t ast.Type) bool { return t.Kind == ast.KindStruct },
		"IsList":      func(t ast.Type) bool { return t.IsList },
		"Ceil":        func(n int) int { return int(math.Ceil(float64(n) / 8.0)) },
	}
	return g
}

func (g *GoGenerator) getGoRpcType(t ast.Type) string {
	if t.Name == "nil" { return "" }
	if t.IsList {
		return util.PascalCase(t.Name) + "List"
	}
	if t.Kind == ast.KindBase {
		return util.PascalCase(t.Name)
	}
	if t.Kind == ast.KindEnum {
		return "U8"
	}
	return util.PascalCase(t.Name)
}

func (g *GoGenerator) getGoLogicType(t ast.Type) string {
	prefix := ""
	if t.IsList { prefix = "[]" }
	
	if t.Kind == ast.KindBase {
		switch t.Name {
		case "i8": return prefix + "int8"
		case "u8": return prefix + "uint8"
		case "i16": return prefix + "int16"
		case "u16": return prefix + "uint16"
		case "i32": return prefix + "int32"
		case "u32": return prefix + "uint32"
		case "i64": return prefix + "int64"
		case "u64": return prefix + "uint64"
		case "f32": return prefix + "float32"
		case "f64": return prefix + "float64"
		case "bool": return prefix + "bool"
		case "text": return prefix + "string"
		case "bin": return prefix + "[]byte"
		}
	}
	
	name := util.PascalCase(t.Name)
	if t.Kind == ast.KindStruct {
		return prefix + "*" + name
	}
	return prefix + name
}

func (g *GoGenerator) getGoType(t ast.Type) string {
	prefix := ""
	if t.IsList { prefix = "[]" }
	
	switch t.Name {
	case "i8": return prefix + "int8"
	case "u8": return prefix + "uint8"
	case "i16": return prefix + "int16"
	case "u16": return prefix + "uint16"
	case "i32": return prefix + "int32"
	case "u32": return prefix + "uint32"
	case "i64": return prefix + "int64"
	case "u64": return prefix + "uint64"
	case "f32": return prefix + "float32"
	case "f64": return prefix + "float64"
	case "bool": return prefix + "bool"
	case "text": return prefix + "string"
	case "bin": return prefix + "[]byte"
	}
	
	if t.Kind == ast.KindStruct {
		return prefix + "*" + util.PascalCase(t.Name)
	}
	return prefix + util.PascalCase(t.Name)
}

func (g *GoGenerator) getGoValue(name string) string {
	switch name {
	case "text": return "\"\""
	case "bin", "nil": return "nil"
	case "bool": return "false"
	case "f32", "f64": return "0.0"
	default: return "0"
	}
}

func (g *GoGenerator) getGoTag(field ast.StructField) string {
	tagKeys := strings.Split(g.Config.GoTag, ",")
	if g.Config.GoTag == "" { return "" }
	
	val := field.Tag
	if val == "" { val = util.SnakeCase(field.Name) }
	
	var res []string
	for _, k := range tagKeys {
		res = append(res, fmt.Sprintf("%s:\"%s\"", strings.TrimSpace(k), val))
	}
	return "`" + strings.Join(res, " ") + "`"
}

type baseTypeInfo struct {
	Name, Go string
	IsFloat  bool
	Eps      string
}

func (g *GoGenerator) Generate(schema *ast.Schema) error {
	targetDir := filepath.Join(g.Config.GoDir, "sb")
	os.MkdirAll(targetDir, 0755)

	pkgName := "sb"

	// 0. Validation
	for _, s := range schema.Structs {
		if len(s.Fields) > 255 {
			return fmt.Errorf("struct %s has %d fields, exceeds limit of 255", s.Name, len(s.Fields))
		}
	}

	// 1. Generate type.go
	types := []baseTypeInfo{
		{"I8", "int8", false, ""}, {"U8", "uint8", false, ""},
		{"I16", "int16", false, ""}, {"U16", "uint16", false, ""},
		{"I32", "int32", false, ""}, {"U32", "uint32", false, ""},
		{"I64", "int64", false, ""}, {"U64", "uint64", false, ""},
		{"F32", "float32", true, "1e-6"}, {"F64", "float64", true, "1e-9"},
		{"Bin", "[]byte", false, ""}, {"Text", "string", false, ""},
	}
	if err := g.executeTemplate("_tpl/type.go", filepath.Join(targetDir, "type.go"), map[string]any{
		"Types":   types,
		"Package": pkgName,
	}); err != nil {
		return err
	}

	// 2. Generate Enums
	if err := g.executeTemplate("_tpl/go.enum.tpl", filepath.Join(targetDir, "enum.go"), map[string]any{
		"Enums":   schema.Enums,
		"Package": pkgName,
	}); err != nil {
		return err
	}

	// 3. Generate Structs
	for _, s := range schema.Structs {
		path := filepath.Join(targetDir, "struct_"+util.SnakeCase(s.Name)+".go")
		if err := g.executeTemplate("_tpl/go.struct.tpl", path, map[string]any{
			"Name":    s.Name,
			"Fields":  s.Fields,
			"Note":    s.Note,
			"Package": pkgName,
		}); err != nil {
			return err
		}
	}

	// 4. Generate API & RPC
	if len(schema.Apis) > 0 {
		modName := g.getModuleName()
		
		// Logic handlers
		for _, api := range schema.Apis {
			filename := "api." + api.Name + ".go"
			logicPath := filepath.Join(targetDir, filename)
			if _, err := os.Stat(logicPath); os.IsNotExist(err) {
				if err := g.executeTemplate("_tpl/go.api.tpl", logicPath, map[string]any{
					"Api":     api,
					"Mod":     modName,
					"Package": pkgName,
				}); err != nil {
					return err
				}
			}
		}

		// API registration
		groups := make(map[string][]ast.Api)
		for _, api := range schema.Apis {
			module := "api"
			if parts := strings.Split(api.Name, "."); len(parts) > 1 {
				module = parts[0]
			}
			groups[module] = append(groups[module], api)
		}
		if err := g.executeTemplate("_tpl/go.api._.tpl", filepath.Join(targetDir, "api._.go"), map[string]any{
			"Apis":    schema.Apis,
			"Groups":  groups,
			"Mod":     modName,
			"Package": pkgName,
		}); err != nil {
			return err
		}

		// RPC Client/Server
		if err := g.executeTemplate("_tpl/go.rpc.tpl", filepath.Join(targetDir, "rpc.go"), map[string]any{
			"Apis":    schema.Apis,
			"Mod":     modName,
			"Package": pkgName,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (g *GoGenerator) getModuleName() string {
	path := filepath.Join(g.Config.GoDir, "go.mod")
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	return ""
}

func (g *GoGenerator) executeTemplate(tplPath, destPath string, data any) error {
	tplContent, err := g.Config.TplFS.ReadFile(tplPath)
	if err != nil { return fmt.Errorf("read embedded template %s: %w", tplPath, err) }
	
	tpl, err := template.New(filepath.Base(tplPath)).Funcs(g.FuncMap).Parse(string(tplContent))
	if err != nil { return err }
	
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil { return err }
	
	return os.WriteFile(destPath, buf.Bytes(), 0644)
}
