package internal

import (
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
	path := filepath.Join(baseDir, "sb", "DOC.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
