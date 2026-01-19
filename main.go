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
	goDir := flag.String("go", "./go", "Output directory for Go code")
	tsDir := flag.String("ts", "./ts", "Output directory for TypeScript code")
	tags := flag.String("tag", "", "Go struct tags (e.g., bson,json)")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: sb [flags] <file.sb>")
		flag.Usage()
		return
	}

	content, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.New(string(content))
	p := parser.New(l)
	schema, err := p.ParseSchema()
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		os.Exit(1)
	}

	cfg := generator.Config{
		GoDir: *goDir,
		TsDir: *tsDir,
		GoTag: *tags,
		TplFS: generator.TplFS,
	}

	goGen := generator.NewGoGenerator(cfg)
	if err := goGen.Generate(schema); err != nil {
		fmt.Printf("Error generating Go code: %v\n", err)
		os.Exit(1)
	}

	tsGen := generator.NewTsGenerator(cfg)
	if err := tsGen.Generate(schema); err != nil {
		fmt.Printf("Error generating TypeScript code: %v\n", err)
		os.Exit(1)
	}

	// Generate Documentation
	if err := generateDoc(schema, cfg); err != nil {
		fmt.Printf("Error generating documentation: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated code.")
}

func generateDoc(schema *ast.Schema, cfg generator.Config) error {
	tplContent, err := cfg.TplFS.ReadFile("_tpl/doc.md.tpl")
	if err != nil {
		return err
	}

	groups := make(map[string][]ast.Api)
	for _, api := range schema.Apis {
		module := "api"
		if parts := strings.Split(api.Name, "."); len(parts) > 1 {
			module = parts[0]
		}
		groups[module] = append(groups[module], api)
	}

	dummyCfg := generator.Config{TplFS: cfg.TplFS}
	goGen := generator.NewGoGenerator(dummyCfg)
	tsGen := generator.NewTsGenerator(dummyCfg)

	funcMap := template.FuncMap{
		"SnakeCase":  util.SnakeCase,
		"PascalCase": util.PascalCase,
		"CamelCase":  util.CamelCase,
		"GoValue":    goGen.FuncMap["GoValue"],
		"TsValue":    tsGen.FuncMap["TsValue"],
	}

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

	os.WriteFile(filepath.Join(cfg.GoDir, "sb", "DOC.md"), buf.Bytes(), 0644)
	os.WriteFile(filepath.Join(cfg.TsDir, "sb", "DOC.md"), buf.Bytes(), 0644)
	return nil
}