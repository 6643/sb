package internal

func (g *TsGenerator) renderStructFile(st TplStruct) string {
	name := PascalCase(st.Name)
	headerWidthsName := g.structHeaderWidthsName(st)
	metaName := CamelCase(st.Name) + "Meta"

	var w sourceWriter
	w.Line("import * as rt from \"./runtime_core\"")
	w.Line("import * as rm from \"./runtime_meta\"")
	for _, line := range g.tsStructImportLines(st) {
		w.Line(line)
	}
	w.Blank()
	w.WriteLineComment("// ", st.Note)
	w.Linef("export interface %s extends rt.Serializable, rt.Deserializable {", name)
	for _, field := range st.Fields {
		w.WriteLineComment("    // ", field.Note)
		w.Linef("    %s: %s;", CamelCase(field.Name), g.structFieldType(field.Type))
	}
	w.Line("}")
	w.Blank()
	w.Linef("const %s = %s;", headerWidthsName, g.structHeaderWidthsLiteral(st))
	w.Blank()
	w.Linef("const %s = rm.defineStruct<%s>({", metaName, name)
	w.Linef("    name: %q,", name)
	w.Linef("    headerWidths: %s,", headerWidthsName)
	w.Line("    create: () => ({")
	for _, field := range st.Fields {
		w.Linef("        %s: %s,", CamelCase(field.Name), g.structDefaultValue(field.Type))
	}
	w.Linef("    }) as any as %s,", name)
	w.Line("    fields: [")
	w.Blank()
	for _, field := range st.Fields {
		w.Linef("        %s,", g.tsFieldMetaExpr(name, field))
	}
	w.Line("    ],")
	w.Line("});")
	w.Blank()
	w.Linef("export const new%s = (): %s => rm.newStruct(%s, get%s, set%s);", name, name, metaName, name, name)
	w.Linef("export const isZero%s = (s: %s | null | undefined): boolean => rm.isZeroStruct(%s, s);", name, name, metaName)
	w.Linef("export const validate%s = (s: %s | null | undefined): rt.Err => rm.validateStruct(%s, s);", name, name, metaName)
	w.Linef("export const get%s = (buf: rt.Buffer): [%s, rt.Err] => rm.getStruct(%s, buf);", name, name, metaName)
	w.Linef("export const set%s = (buf: rt.Buffer, s: %s): rt.Err => rm.setStruct(%s, buf, s);", name, name, metaName)
	w.Linef("export const read%s = (buf: rt.Buffer): [%s, rt.Err] => get%s(buf);", name, name, name)
	w.Linef("export const eq%s = (a: %s | null | undefined, b: %s | null | undefined): boolean => rm.eqStruct(%s, a, b);", name, name, name, metaName)
	w.Blank()
	g.renderStructListBodyHelpers(&w, st)
	return w.String()
}

func (g *TsGenerator) renderStructListBodyHelpers(w *sourceWriter, st TplStruct) {
	name := PascalCase(st.Name)
	w.Linef("export const get%sListBody = (buf: rt.Buffer, state: number): [%s[], rt.Err] => {", name, name)
	w.Linef("    const [list, err] = rt.getDefaultList<%s>(", name)
	w.Line("        buf,")
	w.Line("        state,")
	w.Linef("        () => new%s(),", name)
	w.Linef("        (buf) => read%s(buf),", name)
	w.Line("    );")
	w.Line("    return [list, err];")
	w.Line("};")
	w.Blank()
	w.Linef("export const set%sListBody = (buf: rt.Buffer, state: number, v: %s[]): rt.Err => {", name, name)
	w.Linef("    return rt.setDefaultList<%s>(", name)
	w.Line("        buf,")
	w.Line("        state,")
	w.Line("        v,")
	w.Linef("        (item) => isZero%s(item),", name)
	w.Linef("        (buf, item) => set%s(buf, item),", name)
	w.Line("    );")
	w.Line("}")
}
