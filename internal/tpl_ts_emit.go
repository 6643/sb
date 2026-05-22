package internal

import (
	"os"
	"path/filepath"
	"strings"
)

func (g *TsGenerator) Generate(schema *TplSchema) error {
	targetDir := filepath.Join(g.Config.TsDir, "sb")
	dir, err := newGeneratedDir(targetDir)
	if err != nil {
		return err
	}
	if err := g.removeLegacyAPIWrapperStructFiles(targetDir); err != nil {
		return err
	}
	if err := removeStaleTSGeneratedFiles(targetDir, schema); err != nil {
		return err
	}

	files := g.tsGeneratedFiles(schema)
	if len(schema.Enums) > 0 {
		files = append(files, generatedFile{RelativePath: "enum_smoke.test.ts", Data: []byte(g.renderEnumSmokeTest(schema.Enums)), Perm: 0644})
	}
	if err := removeLegacyTSRuntimeFiles(targetDir, tsAllowedRuntimeFiles(files)); err != nil {
		return err
	}
	if err := dir.WriteAll(files...); err != nil {
		return err
	}

	structFiles := make([]string, 0, len(schema.Structs))
	for _, st := range schema.Structs {
		filename := "struct_" + SnakeCase(st.Name) + ".ts"
		structFiles = append(structFiles, strings.TrimSuffix(filename, ".ts"))
		if err := dir.Write(filename, []byte(g.renderStructFile(st)), 0644); err != nil {
			return err
		}
	}

	allFiles := append([]string{"enum"}, structFiles...)
	if err := dir.Write("_.ts", []byte(g.renderIndexFile(allFiles)), 0644); err != nil {
		return err
	}

	if len(schema.Apis) == 0 {
		return nil
	}
	return dir.WriteAll(
		generatedFile{RelativePath: "rpc.ts", Data: []byte(g.renderRPCFile(schema.Apis)), Perm: 0644},
		generatedFile{RelativePath: "rpc_smoke.test.ts", Data: []byte(g.renderSmokeTest(schema.Apis)), Perm: 0644},
	)
}

func (g *TsGenerator) removeLegacyAPIWrapperStructFiles(targetDir string) error {
	patterns := []string{
		filepath.Join(targetDir, "struct_api_*_req.ts"),
		filepath.Join(targetDir, "struct_api_*_resp.ts"),
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

func removeStaleTSGeneratedFiles(targetDir string, schema *TplSchema) error {
	activeStructs := make(map[string]struct{}, len(schema.Structs))
	for _, st := range schema.Structs {
		activeStructs["struct_"+SnakeCase(st.Name)+".ts"] = struct{}{}
	}
	if err := removeStaleTSStructFiles(targetDir, activeStructs); err != nil {
		return err
	}

	if len(schema.Apis) == 0 {
		if err := removeTSFiles(targetDir, "rpc.ts", "rpc_smoke.test.ts"); err != nil {
			return err
		}
	}
	if len(schema.Enums) == 0 {
		if err := removeTSFiles(targetDir, "enum_smoke.test.ts"); err != nil {
			return err
		}
	}
	return nil
}

func removeStaleTSStructFiles(targetDir string, active map[string]struct{}) error {
	matches, err := filepath.Glob(filepath.Join(targetDir, "struct_*.ts"))
	if err != nil {
		return err
	}
	for _, match := range matches {
		name := filepath.Base(match)
		if strings.HasSuffix(name, ".test.ts") {
			continue
		}
		if _, ok := active[name]; ok {
			continue
		}
		if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func removeTSFiles(targetDir string, names ...string) error {
	for _, name := range names {
		if err := os.Remove(filepath.Join(targetDir, name)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (g *TsGenerator) tsGeneratedFiles(schema *TplSchema) []generatedFile {
	return []generatedFile{
		{RelativePath: "type.ts", Data: []byte(renderTsRuntimeSource()), Perm: 0644},
		{RelativePath: "runtime_core.ts", Data: []byte(renderTsRuntimeCoreSource()), Perm: 0644},
		{RelativePath: "runtime_base.ts", Data: []byte(renderTsRuntimeBaseSource()), Perm: 0644},
		{RelativePath: "runtime_header.ts", Data: []byte(renderTsRuntimeHeaderSource()), Perm: 0644},
		{RelativePath: "runtime_text.ts", Data: []byte(renderTsRuntimeTextSource()), Perm: 0644},
		{RelativePath: "runtime_bin.ts", Data: []byte(renderTsRuntimeBinSource()), Perm: 0644},
		{RelativePath: "runtime_list.ts", Data: []byte(renderTsRuntimeListSource()), Perm: 0644},
		{RelativePath: "runtime_meta.ts", Data: []byte(renderTsRuntimeMetaSource()), Perm: 0644},
		{RelativePath: "runtime_enum.ts", Data: []byte(renderTsRuntimeEnumSource()), Perm: 0644},
		{RelativePath: "runtime_struct.ts", Data: []byte(renderTsRuntimeStructSource()), Perm: 0644},
		{RelativePath: "enum.ts", Data: []byte(g.renderEnumFile(schema.Enums)), Perm: 0644},
	}
}

func tsAllowedRuntimeFiles(files []generatedFile) map[string]struct{} {
	allowed := make(map[string]struct{}, len(files))
	for _, file := range files {
		name := filepath.Base(file.RelativePath)
		if !strings.HasPrefix(name, "runtime_") || !strings.HasSuffix(name, ".ts") {
			continue
		}
		allowed[name] = struct{}{}
	}
	return allowed
}

func removeLegacyTSRuntimeFiles(targetDir string, allowed map[string]struct{}) error {
	matches, err := filepath.Glob(filepath.Join(targetDir, "runtime_*.ts"))
	if err != nil {
		return err
	}
	for _, match := range matches {
		name := filepath.Base(match)
		if strings.HasSuffix(name, ".test.ts") {
			continue
		}
		if _, ok := allowed[name]; ok {
			continue
		}
		if err := os.Remove(match); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
