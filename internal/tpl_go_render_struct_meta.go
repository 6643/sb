package internal

func (g *GoGenerator) renderStructTypeAndMeta(w *sourceWriter, st TplStruct, structName, resetFunc, defaultFunc, isZeroFunc, validateFunc, headerWidthsName string) {
	w.WriteLineCommentWithHead("// ", structName+" ", st.Note)
	w.Linef("type %s struct {", structName)
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		tag := g.getGoTag(field)
		fieldType := g.getGoLogicType(field.Type)
		w.WriteLineComment("\t// ", field.Note)
		if tag == "" {
			w.Linef("\t%s %s", fieldName, fieldType)
			continue
		}
		w.Linef("\t%s %s %s", fieldName, fieldType, tag)
	}
	w.Line("}")
	w.Blank()
	w.Linef("func %s(s *%s) {", resetFunc, structName)
	w.Line("\tif s == nil { return }")
	w.Linef("\t*s = %s{}", structName)
	for _, field := range st.Fields {
		if field.Type.Kind == TplKindEnum && !field.Type.IsList {
			fieldName := PascalCase(field.Name)
			enumName := PascalCase(field.Type.Name)
			w.Linef("\ts.%s = %s()", fieldName, g.enumDefaultName(enumName))
		}
	}
	w.Line("}")
	w.Blank()
	w.Linef("func %s() *%s {", defaultFunc, structName)
	w.Linef("\ts := &%s{}", structName)
	w.Linef("\t%s(s)", resetFunc)
	w.Line("\treturn s")
	w.Line("}")
	w.Blank()
	w.Linef("var %s = [...]uint8{%s}", headerWidthsName, g.structHeaderWidthsLiteral(st))
	w.Blank()
	w.Linef("var %s = DefineStruct[%s](", CamelCase(st.Name)+"Meta", structName)
	for _, field := range st.Fields {
		w.Linef("\t%s,", g.goStructFieldMetaExpr(structName, field))
	}
	w.Line(")")
	w.Blank()
	w.Linef("func %s(s *%s) bool { return IsZeroStruct(%sMeta, s) }", isZeroFunc, structName, CamelCase(st.Name))
	w.Blank()
	w.Linef("func %s(s *%s) error { return ValidateStruct(%sMeta, s) }", validateFunc, structName, CamelCase(st.Name))
	w.Blank()
}
