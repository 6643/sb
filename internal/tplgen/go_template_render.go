package tplgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func (g *GoGenerator) executeTemplate(tplPath, destPath string, data any) error {
	rendered, err := g.renderTemplate(tplPath, data)
	if err != nil {
		return err
	}
	return os.WriteFile(destPath, rendered, 0644)
}

func (g *GoGenerator) renderTemplate(tplPath string, data any) ([]byte, error) {
	tpl, ok := g.tplCache[tplPath]
	if !ok {
		tplContent, err := tplFS.ReadFile(tplPath)
		if err != nil {
			return nil, fmt.Errorf("read embedded template %s: %w", tplPath, err)
		}

		tpl, err = template.New(filepath.Base(tplPath)).Funcs(g.FuncMap).Parse(string(tplContent))
		if err != nil {
			return nil, err
		}
		g.tplCache[tplPath] = tpl
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
