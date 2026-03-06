package internal

import (
	"fmt"
	"path/filepath"
)

func (g *GoGenerator) Generate(schema *TplSchema) error {
	targetDir := filepath.Join(g.Config.GoDir, "sb")
	dir, err := newGeneratedDir(targetDir)
	if err != nil {
		return err
	}

	for _, s := range schema.Structs {
		if len(s.Fields) > 255 {
			return fmt.Errorf("结构体 %s 拥有 %d 个字段，超过限制 (255)", s.Name, len(s.Fields))
		}
	}

	if err := dir.WriteAll(
		generatedFile{RelativePath: "type.go", Data: []byte(renderGoRuntimeSource()), Perm: 0644},
		generatedFile{RelativePath: "enum.go", Data: []byte(g.renderGoEnumFile(schema.Enums)), Perm: 0644},
	); err != nil {
		return err
	}
	for _, s := range schema.Structs {
		if err := dir.Write("struct_"+SnakeCase(s.Name)+".go", []byte(g.renderGoStructFile(s)), 0644); err != nil {
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
		if err := dir.SyncFingerprint(filename, g.renderGoAPIStub(api), 0644); err != nil {
			return err
		}
	}
	if err := dir.WriteAll(
		generatedFile{RelativePath: "api._.go", Data: []byte(g.renderGoAPIHandlers(schema.Apis)), Perm: 0644},
		generatedFile{RelativePath: "rpc.go", Data: []byte(g.renderGoRPCFile(schema.Apis)), Perm: 0644},
	); err != nil {
		return err
	}
	return nil
}
