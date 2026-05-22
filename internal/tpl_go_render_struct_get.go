package internal

func (g *GoGenerator) renderStructGet(w *sourceWriter, st TplStruct, structName, getFunc, resetFunc string, headerBits int, headerWidthsName string, headerFieldCount int) {
	w.Linef("func %s(buf *bytes.Buffer, s *%s) error {", getFunc, structName)
	w.Linef("\tif s == nil { return fmt.Errorf(\"Get%s: nil value\") }", structName)
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		ref := "s." + fieldName
		switch {
		case field.Type.Kind == TplKindStruct && !field.Type.IsList:
			w.Linef("\treuse%s := %s", fieldName, ref)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
			w.Linef("\treuse%s := %s", fieldName, ref)
		case field.Type.IsList:
			w.Linef("\treuse%s := %s", fieldName, ref)
		}
	}
	w.Linef("\t%s(s)", resetFunc)
	w.Linef("\tconst headerBits = %d", headerBits)
	w.Line("\theaderSize := HeaderSize(headerBits)")
	w.Linef("\tif buf.Len() < headerSize { return fmt.Errorf(\"Get%s header: %%d - %%d\", buf.Len(), headerSize) }", structName)
	w.Line("\theader := buf.Next(headerSize)")
	w.Linef("\tvar headerStates [%d]uint8", headerFieldCount)
	w.Linef("\tif err := ReadHeader(header, %s[:], headerStates[:], \"%s header\"); err != nil { return err }", headerWidthsName, structName)
	for i, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		switch g.tagWidth(field.Type) {
		case 1:
			if field.Type.Name == "bool" {
				w.Linef("\t%sState := headerStates[%d] != 0", CamelCase(fieldName), i)
			} else {
				w.Linef("\t%sPresent := headerStates[%d] != 0", CamelCase(fieldName), i)
			}
		case 2:
			w.Linef("\t%sState := headerStates[%d]", CamelCase(fieldName), i)
		}
	}
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		ref := "s." + fieldName
		switch {
		case field.Type.Name == "bool":
			w.Linef("\t%s = %sState", ref, CamelCase(fieldName))
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "text":
			w.Linef("\tvalue%s, err := GetText(buf, %sState)", fieldName, CamelCase(fieldName))
			w.Linef("\tif err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\t%s = value%s", ref, fieldName)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
			w.Linef("\tvalue%s, err := GetBinInto(buf, %sState, reuse%s)", fieldName, CamelCase(fieldName), fieldName)
			w.Linef("\tif err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\t%s = value%s", ref, fieldName)
		case field.Type.Kind == TplKindBase && field.Type.IsList:
			w.Linef("\tvalue%s, err := %s", fieldName, g.listGetReuseExpr(field.Type, CamelCase(fieldName)+"State", "reuse"+fieldName))
			w.Linef("\tif err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\t%s = value%s", ref, fieldName)
		case field.Type.Kind == TplKindStruct && field.Type.IsList:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tvalue%s, err := get%sListBodyReuse(buf, %sState, reuse%s)", fieldName, typeName, CamelCase(fieldName), fieldName)
			w.Linef("\tif err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\t%s = value%s", ref, fieldName)
		case field.Type.Kind == TplKindEnum && field.Type.IsList:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tvalue%s, err := get%sListBodyReuse(buf, %sState, reuse%s)", fieldName, typeName, CamelCase(fieldName), fieldName)
			w.Linef("\tif err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\t%s = value%s", ref, fieldName)
		case field.Type.Kind == TplKindStruct:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tif %sPresent {", CamelCase(fieldName))
			w.Linef("\t\tvalue%s := reuse%s", fieldName, fieldName)
			w.Linef("\t\tif value%s == nil { value%s = %s() }", fieldName, fieldName, g.structDefaultName(typeName))
			w.Linef("\t\tif err := %s(buf, value%s); err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", g.structGetName(field.Type.Name), fieldName, structName, fieldName)
			w.Linef("\t\t%s = value%s", ref, fieldName)
			w.Line("\t}")
		case field.Type.Kind == TplKindEnum:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tif %sPresent {", CamelCase(fieldName))
			w.Line("\t\tvalue, err := GetU8(buf)")
			w.Linef("\t\tif err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\t\titem := %s(value)", typeName)
			w.Linef("\t\tif !%s(item) { return fmt.Errorf(\"Get%s %s: 非法枚举值: %%d\", item) }", g.enumIsName(typeName), structName, fieldName)
			w.Linef("\t\t%s = item", ref)
			w.Line("\t}")
		default:
			if width, getter := g.primitiveGetter(field.Type.Name); width > 0 {
				w.Linef("\tif %sPresent {", CamelCase(fieldName))
				w.Linef("\t\tvalue, err := %s(buf)", getter)
				w.Linef("\t\tif err != nil { return fmt.Errorf(\"Get%s %s: %%w\", err) }", structName, fieldName)
				w.Linef("\t\t%s = value", ref)
				w.Line("\t}")
			}
		}
	}
	w.Line("\treturn nil")
	w.Line("}")
	w.Blank()
}
