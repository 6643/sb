package internal

import (
	"fmt"
)

// Generate 使用字符串生成后端生成 Go/TS 代码与文档。
func Generate(schema *IRSchema, cfg Config) error {
	tplSchema := fromIRSchema(schema)

	goGen := NewGoGenerator(cfg)
	if err := goGen.Generate(tplSchema); err != nil {
		return fmt.Errorf("go generation: %w", err)
	}

	if cfg.GoV2 {
		goV2Gen := NewGoV2Generator(cfg)
		if err := goV2Gen.Generate(tplSchema); err != nil {
			return fmt.Errorf("go v2 generation: %w", err)
		}
	}

	tsGen := NewTsGenerator(cfg)
	if err := tsGen.Generate(tplSchema); err != nil {
		return fmt.Errorf("ts generation: %w", err)
	}

	if cfg.TsV2 {
		tsV2Gen := NewTsV2Generator(cfg)
		if err := tsV2Gen.Generate(tplSchema); err != nil {
			return fmt.Errorf("ts v2 generation: %w", err)
		}
	}

	if err := generateDoc(tplSchema, cfg); err != nil {
		return fmt.Errorf("doc generation: %w", err)
	}
	return nil
}
