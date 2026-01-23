package generator

import (
	"sb/internal/ast"
	"embed"
)

// Config 代码生成配置
type Config struct {
	GoDir string   // Go 代码输出目录
	TsDir string   // TypeScript 代码输出目录
	GoTag string   // 附加的 Go struct tag (如 "bson,json")
	TplFS embed.FS // 嵌入的模板文件系统
}

// Generator 代码生成器接口
type Generator interface {
	Generate(schema *ast.Schema) error
}
