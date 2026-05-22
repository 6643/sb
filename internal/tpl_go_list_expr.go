package internal

import "fmt"

func (g *GoGenerator) listSizeExpr(t TplType, ref string) string {
	switch t.Name {
	case "bool":
		return fmt.Sprintf("SizeBoolList(%s)", ref)
	case "text":
		return fmt.Sprintf("SizeTextList(%s)", ref)
	case "bin":
		return fmt.Sprintf("SizeBinList(%s)", ref)
	default:
		if t.Kind == TplKindBase {
			width, _ := goBaseEncodedWidth(t.Name)
			return fmt.Sprintf("SizeZeroList(%s, %d)", ref, width)
		}
		defaultFactory, callbacks := g.bitmapListCallbacks(t)
		_ = defaultFactory
		return fmt.Sprintf("SizeDefaultList(%s, %s, %s)", ref, callbacks[0], callbacks[1])
	}
}

func (g *GoGenerator) listGetExpr(t TplType, stateVar string) string {
	return g.listGetReuseExpr(t, stateVar, "nil")
}

func (g *GoGenerator) listGetReuseExpr(t TplType, stateVar, dstVar string) string {
	switch t.Name {
	case "bool":
		return fmt.Sprintf("GetBoolListInto(buf, %s, %s)", stateVar, dstVar)
	case "text":
		return fmt.Sprintf("GetTextListInto(buf, %s, %s)", stateVar, dstVar)
	case "bin":
		return fmt.Sprintf("GetBinListInto(buf, %s, %s)", stateVar, dstVar)
	default:
		if t.Kind == TplKindBase {
			_, getter := g.primitiveGetter(t.Name)
			return fmt.Sprintf("GetZeroListInto(buf, %s, %s, %s)", stateVar, dstVar, getter)
		}
		defaultFactory, callbacks := g.bitmapListCallbacks(t)
		return fmt.Sprintf("GetDefaultListInto(buf, %s, %s, %s, %s)", stateVar, dstVar, defaultFactory, callbacks[2])
	}
}

func (g *GoGenerator) listSetCall(t TplType, ref, stateVar string) string {
	switch t.Name {
	case "bool":
		return fmt.Sprintf("SetBoolList(buf, %s, %s)", stateVar, ref)
	case "text":
		return fmt.Sprintf("SetTextList(buf, %s, %s)", stateVar, ref)
	case "bin":
		return fmt.Sprintf("SetBinList(buf, %s, %s)", stateVar, ref)
	default:
		if t.Kind == TplKindBase {
			width, _ := goBaseEncodedWidth(t.Name)
			return fmt.Sprintf("SetZeroList(buf, %s, %s, %d, %s)", stateVar, ref, width, g.primitiveSetter(t.Name))
		}
		_, callbacks := g.bitmapListCallbacks(t)
		return fmt.Sprintf("SetDefaultList(buf, %s, %s, %s, %s, %s)", stateVar, ref, callbacks[0], callbacks[1], callbacks[3])
	}
}

func (g *GoGenerator) bitmapListCallbacks(t TplType) (string, [4]string) {
	elemType := g.getGoType(TplType{Name: t.Name, Kind: t.Kind, IsList: false})
	switch t.Kind {
	case TplKindEnum:
		enumName := PascalCase(t.Name)
		return fmt.Sprintf("func() %s { return %s() }", elemType, g.enumDefaultName(enumName)), [4]string{
			fmt.Sprintf("func(item %s) bool { return %s(item) }", elemType, g.enumIsDefaultName(enumName)),
			fmt.Sprintf("func(item %s) (int, error) { if !%s(item) { return 0, fmt.Errorf(\"非法枚举值: %%d\", item) }; return 1, nil }", elemType, g.enumIsAssignableName(enumName)),
			fmt.Sprintf("func(buf *bytes.Buffer) (%s, error) { value, err := GetU8(buf); if err != nil { return 0, err }; item := %s(value); if !Is%s(item) { return 0, fmt.Errorf(\"非法枚举值: %%d\", item) }; return item, nil }", elemType, elemType, enumName),
			fmt.Sprintf("func(buf *bytes.Buffer, item %s) error { if !Is%s(item) { return fmt.Errorf(\"非法枚举值: %%d\", item) }; return SetU8(buf, uint8(item)) }", elemType, enumName),
		}
	case TplKindStruct:
		typeName := PascalCase(t.Name)
		return "nil", [4]string{
			fmt.Sprintf("func(item %s) bool { return %s(item) }", elemType, g.structIsZeroName(typeName)),
			fmt.Sprintf("func(item %s) (int, error) { return %s(item) }", elemType, g.structSizeName(t.Name)),
			fmt.Sprintf("func(buf *bytes.Buffer) (%s, error) { return %s(buf) }", elemType, g.structReadName(t.Name)),
			fmt.Sprintf("func(buf *bytes.Buffer, item %s) error { return %s(buf, item) }", elemType, g.structSetName(t.Name)),
		}
	default:
		defaultValue := g.getGoValue(t.Name)
		_, getter := g.primitiveGetter(t.Name)
		setter := g.primitiveSetter(t.Name)
		width, _ := goBaseEncodedWidth(t.Name)
		if t.Name == "f32" || t.Name == "f64" {
			defaultValue = "0"
		}
		return fmt.Sprintf("func() %s { return %s }", elemType, defaultValue), [4]string{
			fmt.Sprintf("func(item %s) bool { return item == 0 }", elemType),
			fmt.Sprintf("func(item %s) (int, error) { return %d, nil }", elemType, width),
			fmt.Sprintf("func(buf *bytes.Buffer) (%s, error) { return %s(buf) }", elemType, getter),
			fmt.Sprintf("func(buf *bytes.Buffer, item %s) error { return %s(buf, item) }", elemType, setter),
		}
	}
}

func (g *GoGenerator) eqExpr(t TplType, left, right string) string {
	switch {
	case t.Kind == TplKindStruct && t.IsList:
		return fmt.Sprintf("(slices.EqualFunc(%s, %s, Eq%s))", left, right, PascalCase(t.Name))
	case t.Kind == TplKindStruct:
		return fmt.Sprintf("(Eq%s(%s, %s))", PascalCase(t.Name), left, right)
	case t.Kind == TplKindEnum && t.IsList:
		return fmt.Sprintf("(%s(%s, %s))", g.enumEqListName(PascalCase(t.Name)), left, right)
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("(%s(%s, %s))", g.enumEqValueName(PascalCase(t.Name)), left, right)
	case t.Kind == TplKindBase && t.IsList && t.Name == "bin":
		return fmt.Sprintf("(eqBinList(%s, %s))", left, right)
	case t.Kind == TplKindBase && t.IsList && t.Name == "f32":
		return fmt.Sprintf("(slices.EqualFunc(%s, %s, eqF32))", left, right)
	case t.Kind == TplKindBase && t.IsList && t.Name == "f64":
		return fmt.Sprintf("(slices.EqualFunc(%s, %s, eqF64))", left, right)
	case t.Kind == TplKindBase && t.IsList:
		return fmt.Sprintf("(slices.Equal(%s, %s))", left, right)
	case t.Name == "bin":
		return fmt.Sprintf("(eqBin(%s, %s))", left, right)
	case t.Name == "f32":
		return fmt.Sprintf("(eqF32(%s, %s))", left, right)
	case t.Name == "f64":
		return fmt.Sprintf("(eqF64(%s, %s))", left, right)
	default:
		return fmt.Sprintf("(%s == %s)", left, right)
	}
}

func (g *GoGenerator) goStructFieldMetaExpr(structName string, field TplStructField) string {
	fieldName := PascalCase(field.Name)
	label := fieldName
	ref := fmt.Sprintf("func(s *%s) %s { return s.%s }", structName, g.getGoLogicType(field.Type), fieldName)

	switch {
	case field.Type.Name == "text" && !field.Type.IsList:
		return fmt.Sprintf("TextField[%s](%q, %s)", structName, label, ref)
	case field.Type.Name == "bin" && !field.Type.IsList:
		return fmt.Sprintf("BinField[%s](%q, %s)", structName, label, ref)
	case field.Type.Kind == TplKindEnum && !field.Type.IsList:
		typeName := PascalCase(field.Type.Name)
		return fmt.Sprintf("EnumField[%s, %s](%q, %s, %s, %s, %s)", structName, g.getGoLogicType(field.Type), label, ref, g.enumIsAssignableName(typeName), g.enumIsDefaultName(typeName), g.enumEqValueName(typeName))
	case field.Type.Kind == TplKindStruct && !field.Type.IsList:
		typeName := PascalCase(field.Type.Name)
		return fmt.Sprintf("PtrField[%s, %s](%q, func(s *%s) *%s { return s.%s }, %s, %s, %s)", structName, typeName, label, structName, typeName, fieldName, g.structValidateName(typeName), g.structIsZeroName(typeName), g.structEqName(field.Type.Name))
	case field.Type.IsList:
		return fmt.Sprintf("SliceField[%s, %s](%q, func(s *%s) %s { return s.%s }, %s, %s)", structName, g.goListElemType(field.Type), label, structName, g.getGoLogicType(field.Type), fieldName, g.goListValidateExpr(field.Type), g.goListEqExpr(field.Type))
	default:
		return fmt.Sprintf("ScalarField[%s, %s](%q, %s)", structName, g.getGoLogicType(field.Type), label, ref)
	}
}

func (g *GoGenerator) goListElemType(t TplType) string {
	switch {
	case t.Kind == TplKindStruct:
		return "*" + PascalCase(t.Name)
	case t.Kind == TplKindEnum:
		return PascalCase(t.Name)
	default:
		return g.getGoType(TplType{Name: t.Name, Kind: t.Kind})
	}
}

func (g *GoGenerator) goListValidateExpr(t TplType) string {
	switch {
	case t.Kind == TplKindStruct:
		typeName := PascalCase(t.Name)
		return fmt.Sprintf("func(values []*%s) error { _, err := size%sListBody(values); return err }", typeName, typeName)
	case t.Kind == TplKindEnum:
		typeName := PascalCase(t.Name)
		return fmt.Sprintf("func(values []%s) error { _, err := size%sListBody(values); return err }", typeName, typeName)
	case t.Name == "bool":
		return "func(values []bool) error { _, err := SizeBoolList(values); return err }"
	case t.Name == "text":
		return "func(values []string) error { _, err := SizeTextList(values); return err }"
	case t.Name == "bin":
		return "func(values [][]byte) error { _, err := SizeBinList(values); return err }"
	default:
		width, _ := goBaseEncodedWidth(t.Name)
		return fmt.Sprintf("func(values %s) error { _, err := SizeZeroList(values, %d); return err }", g.getGoLogicType(t), width)
	}
}

func (g *GoGenerator) goListEqExpr(t TplType) string {
	switch {
	case t.Kind == TplKindStruct:
		typeName := PascalCase(t.Name)
		return fmt.Sprintf("func(a, b []*%s) bool { return slices.EqualFunc(a, b, %s) }", typeName, g.structEqName(t.Name))
	case t.Kind == TplKindEnum:
		return g.enumEqListName(PascalCase(t.Name))
	case t.Name == "bin":
		return "eqBinList"
	case t.Name == "f32":
		return "func(a, b []float32) bool { return slices.EqualFunc(a, b, eqF32) }"
	case t.Name == "f64":
		return "func(a, b []float64) bool { return slices.EqualFunc(a, b, eqF64) }"
	default:
		return fmt.Sprintf("func(a, b %s) bool { return slices.Equal(a, b) }", g.getGoLogicType(t))
	}
}
