package internal

import (
	"fmt"
	"strings"
)

func (g *TsGenerator) renderEnumFile(enums []TplEnum) string {
	var w sourceWriter
	w.Line("import * as rt from \"./runtime_core\"")
	w.Line("import * as rm from \"./runtime_meta\"")
	w.Blank()
	for _, enum := range enums {
		enumName := PascalCase(enum.Name)
		w.WriteLineComment("// ", enum.Note)
		w.Linef("export enum %s {", enumName)
		for _, child := range enum.Children {
			w.WriteLineComment("    // ", child.Note)
			w.Linef("    %s = %d,", PascalCase(child.Name), child.ID)
		}
		w.Line("}")
		w.Blank()
		w.Linef("const %sMeta = rm.defineEnum<%s>(%s.%s, [", CamelCase(enum.Name), enumName, enumName, PascalCase(enum.Children[0].Name))
		for _, child := range enum.Children {
			w.Linef("    %s.%s,", enumName, PascalCase(child.Name))
		}
		w.Line("] as const);")
		w.Blank()
		w.Linef("export const Default%s = (): %s => %sMeta.defaultValue;", enumName, enumName, CamelCase(enum.Name))
		w.Linef("export const Is%s = (v: %s): boolean => rm.isEnum(%sMeta, v);", enumName, enumName, CamelCase(enum.Name))
		w.Linef("export const Normalize%s = (v: %s): %s => rm.normalizeEnum(%sMeta, v);", enumName, enumName, enumName, CamelCase(enum.Name))
		w.Linef("export const IsDefault%s = (v: %s): boolean => rm.isDefaultEnum(%sMeta, v);", enumName, enumName, CamelCase(enum.Name))
		w.Linef("export const IsAssignable%s = (v: %s): boolean => rm.isAssignableEnum(%sMeta, v);", enumName, enumName, CamelCase(enum.Name))
		w.Linef("export const eq%sValue = (a: %s, b: %s): boolean => rm.eqEnumValue(%sMeta, a, b);", enumName, enumName, enumName, CamelCase(enum.Name))
		w.Linef("export const eq%sList = (a: %s[], b: %s[]): boolean => rm.eqEnumList(%sMeta, a, b);", enumName, enumName, enumName, CamelCase(enum.Name))
		w.Blank()
		w.Linef("export const get%sListBody = (buf: rt.Buffer, state: number): [%s[], rt.Err] => rm.getEnumList(%sMeta, buf, state);", enumName, enumName, CamelCase(enum.Name))
		w.Blank()
		w.Linef("export const set%sListBody = (buf: rt.Buffer, state: number, v: %s[]): rt.Err => rm.setEnumList(%sMeta, buf, state, v);", enumName, enumName, CamelCase(enum.Name))
		w.Blank()
	}
	return w.String()
}

func (g *TsGenerator) renderIndexFile(files []string) string {
	var w sourceWriter
	w.Line("export * from \"./type\"")
	for _, file := range files {
		w.Linef("export * from \"./%s\"", file)
	}
	return w.String()
}

func (g *TsGenerator) structHeaderWidthsName(st TplStruct) string {
	return CamelCase(st.Name) + "HeaderWidths"
}

func (g *TsGenerator) structHeaderWidthsLiteral(st TplStruct) string {
	parts := make([]string, 0, len(st.Fields))
	for _, field := range st.Fields {
		parts = append(parts, fmt.Sprintf("%d", g.tagWidth(field.Type)))
	}
	return "[" + strings.Join(parts, ", ") + "] as const"
}

func (g *TsGenerator) renderEnumSmokeTest(enums []TplEnum) string {
	var w sourceWriter
	w.Line("import { describe, expect, test } from \"bun:test\";")
	w.Blank()
	w.Line("import * as _ from \"./_\";")
	w.Blank()
	w.Line("describe(\"enum smoke\", () => {")
	for _, enum := range enums {
		if len(enum.Children) == 0 {
			continue
		}
		first := PascalCase(enum.Children[0].Name)
		last := PascalCase(enum.Children[len(enum.Children)-1].Name)
		enumName := PascalCase(enum.Name)
		w.Linef("    test(\"%s validator accepts generated values\", () => {", enumName)
		w.Linef("        expect(_.Is%s(_.%s.%s)).toBe(true);", enumName, enumName, first)
		w.Linef("        expect(_.Is%s(_.%s.%s)).toBe(true);", enumName, enumName, last)
		w.Linef("        expect(_.Is%s(255 as any)).toBe(false);", enumName)
		w.Line("    });")
	}
	w.Line("});")
	return w.String()
}
