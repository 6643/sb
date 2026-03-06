package main

import (
	"flag"
	"fmt"
	"os"

	sbint "sb/internal"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("代码生成成功。")
}

func run() error {
	goDir := flag.String("go", "./go", "Go 代码输出目录")
	tsDir := flag.String("ts", "./ts", "TypeScript 代码输出目录")
	goTag := flag.String("tag", "", "Go 结构体 tag，示例: bson,json")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return fmt.Errorf("缺少输入 .sb 文件")
	}

	schema, err := parseAndResolve(args[0])
	if err != nil {
		return err
	}

	cfg := sbint.Config{GoDir: *goDir, TsDir: *tsDir, GoTag: *goTag}
	if err := sbint.Generate(schema, cfg); err != nil {
		return fmt.Errorf("模板后端生成失败: %w", err)
	}

	return nil
}

func parseAndResolve(schemaPath string) (*sbint.IRSchema, error) {
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("读取 schema 文件失败: %w", err)
	}

	l := sbint.New(string(content))
	p := sbint.NewParser(l)
	astSchema, err := p.ParseSchema()
	if err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	resolved, err := sbint.Resolve(astSchema)
	if err != nil {
		return nil, fmt.Errorf("语义校验失败: %w", err)
	}
	return resolved, nil
}
