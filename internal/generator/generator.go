package generator

import (
	"sb/internal/ast"
	"embed"
)

type Config struct {
	GoDir string
	TsDir string
	GoTag string
	TplFS embed.FS
}

type Generator interface {
	Generate(schema *ast.Schema) error
}
