package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

func groupAPIs(apis []TplApi) map[string][]TplApi {
	groups := make(map[string][]TplApi)
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
	dir, err := newGeneratedDir(filepath.Join(baseDir, "sb"))
	if err != nil {
		return err
	}
	return dir.Write("DOC.md", data, 0644)
}

func writeFileIfChanged(path string, data []byte, perm os.FileMode) error {
	existing, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existing, data) {
		return nil
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(path, data, perm)
}
