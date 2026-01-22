package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sb/internal/ast"
	"sb/internal/generator"
	"sb/internal/lexer"
	"sb/internal/parser"
	"sb/internal/util"
	"strings"
	"text/template"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Successfully generated code.")
}

func run() error {
	goDir := flag.String("go", "./go", "Output directory for Go code")
	tsDir := flag.String("ts", "./ts", "Output directory for TypeScript code")
	tags := flag.String("tag", "", "Go struct tags (e.g., bson,json)")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return fmt.Errorf("missing input file")
	}

	schema, err := parseSchema(args[0])
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	cfg := generator.Config{
		GoDir: *goDir,
		TsDir: *tsDir,
		GoTag: *tags,
		TplFS: generator.TplFS,
	}

	if err := generateCode(schema, cfg); err != nil {
		return err
	}

	return nil
}

func parseSchema(filename string) (*ast.Schema, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	return p.ParseSchema()
}

func generateCode(schema *ast.Schema, cfg generator.Config) error {
	goGen := generator.NewGoGenerator(cfg)
	if err := goGen.Generate(schema); err != nil {
		return fmt.Errorf("go generation: %w", err)
	}

	tsGen := generator.NewTsGenerator(cfg)
	if err := tsGen.Generate(schema); err != nil {
		return fmt.Errorf("ts generation: %w", err)
	}

	if err := generateDoc(schema, cfg); err != nil {
		return fmt.Errorf("documentation: %w", err)
	}

	return nil
}

func generateDoc(schema *ast.Schema, cfg generator.Config) error {
	tplContent, err := cfg.TplFS.ReadFile("_tpl/doc.md.tpl")
	if err != nil {
		return err
	}

	groups := groupApis(schema.Apis)
	funcMap := createFuncMap(cfg)

	tpl, err := template.New("doc").Funcs(funcMap).Parse(string(tplContent))
	if err != nil {
		return err
	}

	data := map[string]any{
		"Apis":    schema.Apis,
		"Enums":   schema.Enums,
		"Structs": schema.Structs,
		"Note":    schema.Note,
		"Groups":  groups,
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return err
	}

	if err := writeDocFile(cfg.GoDir, buf.Bytes()); err != nil {
		return err
	}
	if err := writeDocFile(cfg.TsDir, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func groupApis(apis []ast.Api) map[string][]ast.Api {
	groups := make(map[string][]ast.Api)
	for _, api := range apis {
		module := "api"
		if parts := strings.Split(api.Name, "."); len(parts) > 1 {
			module = parts[0]
		}
		groups[module] = append(groups[module], api)
	}
	return groups
}

func createFuncMap(cfg generator.Config) template.FuncMap {
	dummyCfg := generator.Config{TplFS: cfg.TplFS}
	goGen := generator.NewGoGenerator(dummyCfg)
	tsGen := generator.NewTsGenerator(dummyCfg)

	return template.FuncMap{
		"SnakeCase":  util.SnakeCase,
		"PascalCase": util.PascalCase,
		"CamelCase":  util.CamelCase,
		"GoValue":    goGen.FuncMap["GoValue"],
		"TsValue":    tsGen.FuncMap["TsValue"],
	}
}

func writeDocFile(baseDir string, data []byte) error {
	path := filepath.Join(baseDir, "sb", "DOC.md")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}