package tplgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sb/internal/util"
	"strings"
)

type baseTypeInfo struct {
	Name    string
	Go      string
	IsFloat bool
	Eps     string
}

func (g *GoGenerator) Generate(schema *Schema) error {
	targetDir := filepath.Join(g.Config.GoDir, "sb")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create go target dir %s: %w", targetDir, err)
	}

	for _, s := range schema.Structs {
		if len(s.Fields) > 255 {
			return fmt.Errorf("结构体 %s 拥有 %d 个字段，超过限制 (255)", s.Name, len(s.Fields))
		}
	}

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
		"Package": "sb",
	}); err != nil {
		return err
	}

	if err := g.executeTemplate("_tpl/go.enum.tpl", filepath.Join(targetDir, "enum.go"), map[string]any{
		"Enums":   schema.Enums,
		"Package": "sb",
	}); err != nil {
		return err
	}

	for _, s := range schema.Structs {
		path := filepath.Join(targetDir, "struct_"+util.SnakeCase(s.Name)+".go")
		if err := g.executeTemplate("_tpl/go.struct.tpl", path, map[string]any{
			"Name":    s.Name,
			"Fields":  s.Fields,
			"Note":    s.Note,
			"Package": "sb",
		}); err != nil {
			return err
		}
	}

	if len(schema.Apis) == 0 {
		return nil
	}
	if err := g.migrateAPIFilesToSB(filepath.Join(g.Config.GoDir, "api"), targetDir, schema.Apis); err != nil {
		return err
	}

	for _, api := range schema.Apis {
		filename := "api." + api.Name + ".go"
		logicPath := filepath.Join(targetDir, filename)

		generatedBody, err := g.renderTemplate("_tpl/go.api.tpl", map[string]any{
			"Api":     api,
			"Package": "sb",
		})
		if err != nil {
			return err
		}
		generatedContent := withFingerprint(generatedBody)

		existingContent, err := os.ReadFile(logicPath)
		if os.IsNotExist(err) {
			if err := os.WriteFile(logicPath, generatedContent, 0644); err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return fmt.Errorf("read logic file %s: %w", logicPath, err)
		}

		// 兼容历史无指纹文件：当内容仍等于当前生成结果时，自动升级为带指纹版本。
		if bytes.Equal(existingContent, generatedBody) {
			if err := os.WriteFile(logicPath, generatedContent, 0644); err != nil {
				return err
			}
			continue
		}

		existingBody, recordedHash, ok := splitFingerprint(existingContent)
		if !ok {
			continue
		}
		if hashContent(existingBody) != recordedHash {
			continue
		}
		if bytes.Equal(existingContent, generatedContent) {
			continue
		}
		if err := os.WriteFile(logicPath, generatedContent, 0644); err != nil {
			return err
		}
	}

	groups := make(map[string][]Api)
	for _, api := range schema.Apis {
		module := "api"
		parts := strings.Split(api.Name, ".")
		if len(parts) > 1 {
			module = parts[0]
		}
		groups[module] = append(groups[module], api)
	}

	if err := g.executeTemplate("_tpl/go.api._.tpl", filepath.Join(targetDir, "api._.go"), map[string]any{
		"Apis":    schema.Apis,
		"Groups":  groups,
		"Package": "sb",
	}); err != nil {
		return err
	}

	if err := g.executeTemplate("_tpl/go.rpc.tpl", filepath.Join(targetDir, "rpc.go"), map[string]any{
		"Apis":    schema.Apis,
		"Package": "sb",
	}); err != nil {
		return err
	}

	return nil
}
