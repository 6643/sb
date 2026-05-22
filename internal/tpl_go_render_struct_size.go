package internal

func (g *GoGenerator) renderStructSize(w *sourceWriter, st TplStruct, structName, sizeFunc, sizeValidatedFunc, validateFunc string, headerBits int) {
	w.Linef("func %s(s *%s) (int, error) {", sizeValidatedFunc, structName)
	w.Linef("\tif s == nil { return 0, fmt.Errorf(\"Size%s: nil value\") }", structName)
	w.Linef("\tsize := HeaderSize(%d)", headerBits)
	for _, field := range st.Fields {
		fieldName := PascalCase(field.Name)
		ref := "s." + fieldName
		sizeVar := "fieldSize" + fieldName
		switch {
		case field.Type.Name == "bool":
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "text":
			w.Linef("\t%s, err := SizeText(%s)", sizeVar, ref)
			w.Linef("\tif err != nil { return 0, fmt.Errorf(\"Size%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\tsize += %s", sizeVar)
		case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
			w.Linef("\t%s, err := SizeBin(%s)", sizeVar, ref)
			w.Linef("\tif err != nil { return 0, fmt.Errorf(\"Size%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\tsize += %s", sizeVar)
		case field.Type.Kind == TplKindBase && field.Type.IsList:
			w.Linef("\t%s, err := %s", sizeVar, g.listSizeExpr(field.Type, ref))
			w.Linef("\tif err != nil { return 0, fmt.Errorf(\"Size%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\tsize += %s", sizeVar)
		case field.Type.Kind == TplKindStruct && field.Type.IsList:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\t%s, err := size%sListBody(%s)", sizeVar, typeName, ref)
			w.Linef("\tif err != nil { return 0, fmt.Errorf(\"Size%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\tsize += %s", sizeVar)
		case field.Type.Kind == TplKindEnum && field.Type.IsList:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\t%s, err := size%sListBody(%s)", sizeVar, typeName, ref)
			w.Linef("\tif err != nil { return 0, fmt.Errorf(\"Size%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\tsize += %s", sizeVar)
		case field.Type.Kind == TplKindStruct:
			typeName := PascalCase(field.Type.Name)
			w.Linef("\tif !%s(%s) {", g.structIsZeroName(typeName), ref)
			w.Linef("\t\t%s, err := %s(%s)", sizeVar, g.structSizeName(field.Type.Name), ref)
			w.Linef("\t\tif err != nil { return 0, fmt.Errorf(\"Size%s %s: %%w\", err) }", structName, fieldName)
			w.Linef("\t\tsize += %s", sizeVar)
			w.Line("\t}")
		case field.Type.Kind == TplKindEnum:
			enumName := PascalCase(field.Type.Name)
			w.Linef("\tif !%s(%s) { size += 1 }", g.enumIsDefaultName(enumName), ref)
		default:
			if width, ok := goBaseEncodedWidth(field.Type.Name); ok {
				w.Linef("\tif %s { size += %d }", g.nonDefaultExpr(field.Type, ref), width)
			}
		}
	}
	w.Line("\treturn size, nil")
	w.Line("}")
	w.Blank()
	w.Linef("func %s(s *%s) (int, error) {", sizeFunc, structName)
	w.Linef("\tif err := %s(s); err != nil { return 0, fmt.Errorf(\"Validate%s: %%w\", err) }", validateFunc, structName)
	w.Linef("\treturn %s(s)", sizeValidatedFunc)
	w.Line("}")
	w.Blank()
}
