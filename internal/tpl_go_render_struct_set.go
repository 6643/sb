package internal

func (g *GoGenerator) renderStructSet(w *sourceWriter, st TplStruct, structName, setFunc, setValidatedFunc, validateFunc, sizeValidatedFunc string, headerSize int, headerWidthsName string, headerFieldCount int) {
	w.Linef("func %s(buf *bytes.Buffer, s *%s) error {", setValidatedFunc, structName)
	w.Linef("\tif s == nil { return fmt.Errorf(\"Set%s: nil value\") }", structName)
	for _, field := range st.Fields {
		if g.tagWidth(field.Type) == 2 {
			w.Linef("\t%s", g.stateInitLine(field, "return err"))
		}
	}
	w.Linef("\tsize, err := %s(s)", sizeValidatedFunc)
	w.Line("\tif err != nil { return err }")
	w.Linef("\tvar headerData [%d]byte", headerSize)
	w.Linef("\tvar headerStates [%d]uint8", headerFieldCount)
	for i, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		ref := "s." + fieldName
		switch {
		case field.Type.Name == "bool":
			w.Linef("\tif %s { headerStates[%d] = 1 }", ref, i)
		case g.tagWidth(field.Type) == 1:
			w.Linef("\tif %s { headerStates[%d] = 1 }", g.nonDefaultExpr(field.Type, ref), i)
		case g.tagWidth(field.Type) == 2:
			stateVar := CamelCase(fieldName) + "State"
			w.Linef("\theaderStates[%d] = %s", i, stateVar)
		}
	}
	w.Line("\tbuf.Grow(size)")
	w.Linef("\tif err := WriteHeader(headerData[:], %s[:], headerStates[:]); err != nil { return fmt.Errorf(\"Set%s write header: %%w\", err) }", headerWidthsName, structName)
	w.Linef("\tif _, err := buf.Write(headerData[:]); err != nil { return fmt.Errorf(\"Set%s write header: %%w\", err) }", structName)
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		ref := "s." + fieldName
		switch {
		case field.Type.Name == "bool":
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "text":
			w.Linef("\tif err := SetText(buf, %sState, %s); err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", CamelCase(fieldName), ref, structName, fieldName)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
			w.Linef("\tif err := SetBin(buf, %sState, %s); err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", CamelCase(fieldName), ref, structName, fieldName)
		case field.Type.Kind == TplKindBase && field.Type.IsList:
			w.Linef("\tif err := %s; err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", g.listSetCall(field.Type, ref, CamelCase(fieldName)+"State"), structName, fieldName)
		case field.Type.Kind == TplKindStruct && field.Type.IsList:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tif err := set%sListBody(buf, %sState, %s); err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", typeName, CamelCase(fieldName), ref, structName, fieldName)
		case field.Type.Kind == TplKindEnum && field.Type.IsList:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tif err := set%sListBody(buf, %sState, %s); err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", typeName, CamelCase(fieldName), ref, structName, fieldName)
		case field.Type.Kind == TplKindStruct:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tif !%s(%s) {", g.structIsZeroName(typeName), ref)
			w.Linef("\t\tif err := %s(buf, %s); err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", g.structSetName(field.Type.Name), ref, structName, fieldName)
			w.Line("\t}")
		case field.Type.Kind == TplKindEnum:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tif !%s(%s) {", g.enumIsDefaultName(typeName), ref)
			w.Linef("\t\tif err := SetU8(buf, uint8(%s(%s))); err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", g.enumNormalizeName(typeName), ref, structName, fieldName)
			w.Line("\t}")
		default:
			if _, setter := g.primitiveGetter(field.Type.Name); setter != "" {
				w.Linef("\tif %s {", g.nonDefaultExpr(field.Type, ref))
				w.Linef("\t\tif err := %s(buf, %s); err != nil { return fmt.Errorf(\"Set%s %s: %%w\", err) }", g.primitiveSetter(field.Type.Name), ref, structName, fieldName)
				w.Line("\t}")
			}
		}
	}
	w.Line("\treturn nil")
	w.Line("}")
	w.Blank()
	w.Linef("func %s(buf *bytes.Buffer, s *%s) error {", setFunc, structName)
	w.Linef("\tif s == nil { return fmt.Errorf(\"Set%s: nil value\") }", structName)
	w.Linef("\tif err := %s(s); err != nil { return fmt.Errorf(\"Validate%s: %%w\", err) }", validateFunc, structName)
	w.Linef("\treturn %s(buf, s)", setValidatedFunc)
	w.Line("}")
	w.Blank()
}
