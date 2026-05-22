package internal

import (
	"os"
	"path/filepath"
)

func (g *GoGenerator) Generate(schema *TplSchema) error {
	targetDir := filepath.Join(g.Config.GoDir, "sb")
	dir, err := newGeneratedDir(targetDir)
	if err != nil {
		return err
	}
	if err := removeLegacyGoLayoutFiles(targetDir); err != nil {
		return err
	}
	if err := removeStaleGoAPIFiles(targetDir, schema.Apis); err != nil {
		return err
	}
	if err := dir.Write("runtime.go", []byte(renderGoRuntimeSource()), 0644); err != nil {
		return err
	}
	if err := dir.Write("schema.gen.go", []byte(g.renderSchemaFile(schema)), 0644); err != nil {
		return err
	}
	if len(schema.Apis) == 0 {
		return nil
	}
	for _, api := range schema.Apis {
		filename := "api." + api.Name + ".go"
		if err := dir.SyncFingerprint(filename, g.renderAPIStub(api), 0644); err != nil {
			return err
		}
	}
	return dir.WriteAll(
		generatedFile{RelativePath: "api._.go", Data: []byte(g.renderAPIFile(schema.Apis)), Perm: 0644},
		generatedFile{RelativePath: "rpc.go", Data: []byte(g.renderRPCFile(schema.Apis)), Perm: 0644},
	)
}

func removeStaleGoAPIFiles(targetDir string, apis []TplApi) error {
	active := make(map[string]struct{}, len(apis))
	for _, api := range apis {
		active["api."+api.Name+".go"] = struct{}{}
	}

	matches, err := filepath.Glob(filepath.Join(targetDir, "api.*.go"))
	if err != nil {
		return err
	}
	for _, match := range matches {
		name := filepath.Base(match)
		if name == "api._.go" {
			continue
		}
		if _, ok := active[name]; ok {
			continue
		}
		if err := removeFingerprintManagedFile(match); err != nil {
			return err
		}
	}

	if len(apis) != 0 {
		return nil
	}
	for _, name := range []string{"api._.go", "rpc.go"} {
		if err := os.Remove(filepath.Join(targetDir, name)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func removeFingerprintManagedFile(path string) error {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	body, recordedHash, ok := splitFingerprint(content)
	if !ok || hashContent(body) != recordedHash {
		return nil
	}
	return os.Remove(path)
}

func removeLegacyGoLayoutFiles(targetDir string) error {
	if err := os.RemoveAll(filepath.Join(targetDir, "rt")); err != nil {
		return err
	}
	patterns := []string{
		filepath.Join(targetDir, "type.go"),
		filepath.Join(targetDir, "enum.go"),
		filepath.Join(targetDir, "struct_*.go"),
	}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return err
		}
		for _, match := range matches {
			if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}
	return nil
}
