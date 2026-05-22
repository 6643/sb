package internal

import (
	"fmt"
	"slices"
	"strings"
)

func (g *TsGenerator) tsStructImportLines(st TplStruct) []string {
	enumImports := map[string]map[string]struct{}{}
	structImports := map[string]map[string]struct{}{}

	add := func(group map[string]map[string]struct{}, file string, names ...string) {
		if _, ok := group[file]; !ok {
			group[file] = map[string]struct{}{}
		}
		for _, name := range names {
			group[file][name] = struct{}{}
		}
	}

	for _, field := range st.Fields {
		typeName := PascalCase(field.Type.Name)
		switch field.Type.Kind {
		case TplKindEnum:
			add(enumImports, "./enum", typeName, "Default"+typeName, "Is"+typeName, "IsAssignable"+typeName, "IsDefault"+typeName, "Normalize"+typeName, "get"+typeName+"ListBody", "set"+typeName+"ListBody")
		case TplKindStruct:
			if typeName == PascalCase(st.Name) {
				continue
			}
			file := "./struct_" + SnakeCase(field.Type.Name)
			add(structImports, file, typeName, "new"+typeName, "isZero"+typeName, "validate"+typeName, "read"+typeName, "set"+typeName, "eq"+typeName, "get"+typeName+"ListBody", "set"+typeName+"ListBody")
		}
	}

	lines := make([]string, 0, len(enumImports)+len(structImports))
	if names, ok := enumImports["./enum"]; ok {
		lines = append(lines, fmt.Sprintf("import { %s } from \"./enum\"", joinImportNames(names)))
	}
	structFiles := make([]string, 0, len(structImports))
	for file := range structImports {
		structFiles = append(structFiles, file)
	}
	slices.Sort(structFiles)
	for _, file := range structFiles {
		lines = append(lines, fmt.Sprintf("import { %s } from %q", joinImportNames(structImports[file]), file))
	}
	return lines
}

func (g *TsGenerator) tsRpcImportLines(apis []TplApi) []string {
	enumImports := map[string]struct{}{}
	structImports := map[string]map[string]struct{}{}
	addStruct := func(file string, names ...string) {
		if _, ok := structImports[file]; !ok {
			structImports[file] = map[string]struct{}{}
		}
		for _, name := range names {
			structImports[file][name] = struct{}{}
		}
	}
	visit := func(t TplType) {
		if t.Kind == TplKindEnum {
			typeName := PascalCase(t.Name)
			for _, name := range []string{
				typeName,
				"Default" + typeName,
				"Is" + typeName,
				"IsAssignable" + typeName,
				"Normalize" + typeName,
				"get" + typeName + "ListBody",
				"set" + typeName + "ListBody",
			} {
				enumImports[name] = struct{}{}
			}
			return
		}
		if t.Kind == TplKindStruct {
			typeName := PascalCase(t.Name)
			file := "./struct_" + SnakeCase(t.Name)
			addStruct(file, typeName, "new"+typeName, "read"+typeName, "set"+typeName, "get"+typeName+"ListBody", "set"+typeName+"ListBody")
		}
	}
	for _, api := range apis {
		for _, arg := range api.Args {
			visit(arg.Type)
		}
		visit(api.Result)
	}
	lines := make([]string, 0, len(structImports)+1)
	if len(enumImports) > 0 {
		lines = append(lines, fmt.Sprintf("import { %s } from \"./enum\"", joinImportNames(enumImports)))
	}
	structFiles := make([]string, 0, len(structImports))
	for file := range structImports {
		structFiles = append(structFiles, file)
	}
	slices.Sort(structFiles)
	for _, file := range structFiles {
		lines = append(lines, fmt.Sprintf("import { %s } from %q", joinImportNames(structImports[file]), file))
	}
	return lines
}

func joinImportNames(values map[string]struct{}) string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	slices.Sort(names)
	return strings.Join(names, ", ")
}
