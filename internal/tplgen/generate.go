package tplgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"sb/internal/gen"
	"sb/internal/ir"
	"sb/internal/util"
)

// Generate 使用模板后端生成 Go/TS 代码与文档。
func Generate(schema *ir.Schema, cfg gen.Config) error {
	tplSchema := fromIRSchema(schema)

	goGen := NewGoGenerator(cfg)
	if err := goGen.Generate(tplSchema); err != nil {
		return fmt.Errorf("go template generation: %w", err)
	}

	tsGen := NewTsGenerator(cfg)
	if err := tsGen.Generate(tplSchema); err != nil {
		return fmt.Errorf("ts template generation: %w", err)
	}

	if err := generateDoc(tplSchema, cfg, goGen.FuncMap, tsGen.FuncMap); err != nil {
		return fmt.Errorf("doc template generation: %w", err)
	}

	return nil
}

func generateDoc(schema *Schema, cfg gen.Config, goFuncs, tsFuncs template.FuncMap) error {
	tplContent, err := tplFS.ReadFile("_tpl/doc.md.tpl")
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"SnakeCase":  util.SnakeCase,
		"PascalCase": util.PascalCase,
		"CamelCase":  util.CamelCase,
		"GoValue":    goFuncs["GoValue"],
		"TsValue":    tsFuncs["TsValue"],
	}

	tpl, err := template.New("doc").Funcs(funcMap).Parse(string(tplContent))
	if err != nil {
		return err
	}

	data := map[string]any{
		"Apis":    schema.Apis,
		"Enums":   schema.Enums,
		"Structs": schema.Structs,
		"Note":    schema.Note,
		"Groups":  groupAPIs(schema.Apis),
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return err
	}

	if err := writeDocFile(cfg.GoDir, buf.Bytes()); err != nil {
		return err
	}
	if err := writeDocFile(cfg.TsDir, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func groupAPIs(apis []Api) map[string][]Api {
	groups := make(map[string][]Api)
	for _, api := range apis {
		module := "api"
		parts := strings.Split(api.Name, ".")
		if len(parts) > 1 {
			module = parts[0]
		}
		groups[module] = append(groups[module], api)
	}
	return groups
}

func writeDocFile(baseDir string, data []byte) error {
	path := filepath.Join(baseDir, "sb", "DOC.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
