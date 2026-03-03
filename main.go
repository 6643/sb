package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sb/internal/gen"
	"sb/internal/ir"
	"sb/internal/lexer"
	"sb/internal/parser"
	"sb/internal/semantic"
	"sb/internal/tplgen"
	"sb/internal/util"
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

	cfg := gen.Config{GoDir: *goDir, TsDir: *tsDir, GoTag: *goTag}
	if err := tplgen.Generate(schema, cfg); err != nil {
		return fmt.Errorf("模板后端生成失败: %w", err)
	}

	return nil
}

func parseAndResolve(schemaPath string) (*ir.Schema, error) {
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("读取 schema 文件失败: %w", err)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	astSchema, err := p.ParseSchema()
	if err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	resolved, err := semantic.Resolve(astSchema)
	if err != nil {
		return nil, fmt.Errorf("语义校验失败: %w", err)
	}
	return resolved, nil
}

func writeDocs(schema *ir.Schema, goDir, tsDir string) error {
	doc := buildDoc(schema)
	if err := writeOneDoc(goDir, doc); err != nil {
		return err
	}
	if err := writeOneDoc(tsDir, doc); err != nil {
		return err
	}
	return nil
}

func writeOneDoc(baseDir string, content []byte) error {
	path := filepath.Join(baseDir, "sb", "DOC.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

func buildDoc(schema *ir.Schema) []byte {
	var b bytes.Buffer
	b.WriteString("# SB API Documentation\n\n")
	if strings.TrimSpace(schema.Note) != "" {
		b.WriteString(schema.Note)
		b.WriteString("\n\n")
	}

	b.WriteString("## APIs\n\n")
	b.WriteString("| Name | Args | Result | Note |\n")
	b.WriteString("| :--- | :--- | :--- | :--- |\n")
	for _, api := range schema.APIs {
		args := make([]string, 0, len(api.Args))
		for _, a := range api.Args {
			args = append(args, a.Name+" "+docType(a.Type))
		}
		result := docType(api.Result)
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", util.SnakeCase(api.Name), strings.Join(args, "<br>"), result, api.Note)
	}

	b.WriteString("\n## Enums\n\n")
	for _, e := range schema.Enums {
		fmt.Fprintf(&b, "### %s\n\n", e.Name)
		if strings.TrimSpace(e.Note) != "" {
			fmt.Fprintf(&b, "> %s\n\n", e.Note)
		}
		b.WriteString("| ID | Name | Note |\n")
		b.WriteString("| :--- | :--- | :--- |\n")
		for _, m := range e.Members {
			fmt.Fprintf(&b, "| %d | %s | %s |\n", m.ID, m.Name, m.Note)
		}
		b.WriteString("\n")
	}

	b.WriteString("## Structs\n\n")
	for _, st := range schema.Structs {
		fmt.Fprintf(&b, "### %s\n\n", st.Name)
		if strings.TrimSpace(st.Note) != "" {
			fmt.Fprintf(&b, "> %s\n\n", st.Note)
		}
		b.WriteString("| Field | Type | Note |\n")
		b.WriteString("| :--- | :--- | :--- |\n")
		for _, f := range st.Fields {
			fmt.Fprintf(&b, "| %s | %s | %s |\n", f.Name, docType(f.Type), f.Note)
		}
		b.WriteString("\n")
	}

	return b.Bytes()
}

func docType(t ir.Type) string {
	name := t.Name
	if t.Kind == ir.KindNil {
		name = "nil"
	}
	if t.IsList {
		return "[" + name + "]"
	}
	return name
}
