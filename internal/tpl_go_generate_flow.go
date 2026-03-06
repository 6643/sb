package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

func (g *GoGenerator) Generate(schema *TplSchema) error {
	targetDir := filepath.Join(g.Config.GoDir, "sb")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create go target dir %s: %w", targetDir, err)
	}

	for _, s := range schema.Structs {
		if len(s.Fields) > 255 {
			return fmt.Errorf("结构体 %s 拥有 %d 个字段，超过限制 (255)", s.Name, len(s.Fields))
		}
	}

	if err := os.WriteFile(filepath.Join(targetDir, "type.go"), []byte(renderGoRuntimeSource()), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(targetDir, "enum.go"), []byte(g.renderGoEnumFile(schema.Enums)), 0644); err != nil {
		return err
	}
	for _, s := range schema.Structs {
		path := filepath.Join(targetDir, "struct_"+SnakeCase(s.Name)+".go")
		if err := os.WriteFile(path, []byte(g.renderGoStructFile(s)), 0644); err != nil {
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
		generatedBody := g.renderGoAPIStub(api)
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
	if err := os.WriteFile(filepath.Join(targetDir, "api._.go"), []byte(g.renderGoAPIHandlers(schema.Apis)), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(targetDir, "rpc.go"), []byte(g.renderGoRPCFile(schema.Apis)), 0644); err != nil {
		return err
	}
	return nil
}
