package internal

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type generatedFile struct {
	RelativePath string
	Data         []byte
	Perm         os.FileMode
}

type generatedDir struct {
	root string
}

func newGeneratedDir(root string) (*generatedDir, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("create target dir %s: %w", root, err)
	}
	return &generatedDir{root: root}, nil
}

func (d *generatedDir) path(relativePath string) string {
	return filepath.Join(d.root, relativePath)
}

func (d *generatedDir) Write(relativePath string, data []byte, perm os.FileMode) error {
	path := d.path(relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create target dir %s: %w", filepath.Dir(path), err)
	}
	if err := writeFileIfChanged(path, data, perm); err != nil {
		return fmt.Errorf("write generated file %s: %w", path, err)
	}
	return nil
}

func (d *generatedDir) WriteAll(files ...generatedFile) error {
	for _, file := range files {
		if err := d.Write(file.RelativePath, file.Data, file.Perm); err != nil {
			return err
		}
	}
	return nil
}

func (d *generatedDir) WriteIfAbsent(relativePath string, data []byte, perm os.FileMode) error {
	path := d.path(relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create target dir %s: %w", filepath.Dir(path), err)
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat generated file %s: %w", path, err)
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("write generated file %s: %w", path, err)
	}
	return nil
}

func (d *generatedDir) SyncFingerprint(relativePath string, body []byte, perm os.FileMode) error {
	path := d.path(relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create target dir %s: %w", filepath.Dir(path), err)
	}

	generatedContent := withFingerprint(body)
	existingContent, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return d.writeDirect(path, generatedContent, perm)
	}
	if err != nil {
		return fmt.Errorf("read generated file %s: %w", path, err)
	}
	if bytes.Equal(existingContent, generatedContent) {
		return nil
	}
	if bytes.Equal(existingContent, body) {
		return d.writeDirect(path, generatedContent, perm)
	}

	existingBody, recordedHash, ok := splitFingerprint(existingContent)
	if !ok {
		return nil
	}
	if hashContent(existingBody) != recordedHash {
		return nil
	}
	return d.writeDirect(path, generatedContent, perm)
}

func (d *generatedDir) writeDirect(path string, data []byte, perm os.FileMode) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("write generated file %s: %w", path, err)
	}
	return nil
}
