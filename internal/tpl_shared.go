package internal

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
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

func orderedGroupKeys(groups map[string][]TplApi) []string {
	keys := make([]string, 0, len(groups))
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func goBaseEncodedWidth(name string) (int, bool) {
	switch name {
	case "i8", "u8":
		return 1, true
	case "i16", "u16":
		return 2, true
	case "i32", "u32", "f32":
		return 4, true
	case "i64", "u64", "f64":
		return 8, true
	default:
		return 0, false
	}
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
