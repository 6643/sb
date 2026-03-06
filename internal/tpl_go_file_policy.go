package internal

import (
	"bytes"
	"os"
	"path/filepath"
)

func (g *GoGenerator) migrateAPIFilesToSB(apiDir, targetDir string, apis []TplApi) error {
	if _, err := os.Stat(apiDir); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	for _, api := range apis {
		filename := "api." + api.Name + ".go"
		oldPath := filepath.Join(apiDir, filename)
		newPath := filepath.Join(targetDir, filename)

		if _, err := os.Stat(oldPath); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
			continue
		} else if err != nil {
			return err
		}

		oldContent, err := os.ReadFile(oldPath)
		if err != nil {
			return err
		}
		newContent, err := os.ReadFile(newPath)
		if err != nil {
			return err
		}
		if bytes.Equal(oldContent, newContent) {
			if err := os.Remove(oldPath); err != nil {
				return err
			}
		}
	}

	oldHandlerPath := filepath.Join(apiDir, "api._.go")
	if _, err := os.Stat(oldHandlerPath); err == nil {
		if err := os.Remove(oldHandlerPath); err != nil {
			return err
		}
	}
	_ = os.Remove(apiDir)
	return nil
}
