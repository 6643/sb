package internal

import "fmt"

func (g *TsGenerator) structFieldType(t TplType) string {
	base := g.getTsType(TplType{Name: t.Name, Kind: t.Kind})
	if t.IsList {
		return base + "[]"
	}
	if t.Kind == TplKindStruct {
		return base + " | null"
	}
	return base
}

func (g *TsGenerator) structDefaultValue(t TplType) string {
	if t.IsList {
		return "[]"
	}
	switch t.Kind {
	case TplKindBase:
		return g.getTsValue(t.Name)
	case TplKindEnum:
		return fmt.Sprintf("Default%s()", PascalCase(t.Name))
	default:
		return "null"
	}
}

func (g *TsGenerator) tsFieldMetaExpr(structName string, field TplStructField) string {
	key := CamelCase(field.Name)
	label := PascalCase(field.Name)
	t := field.Type
	typeName := PascalCase(t.Name)
	tsType := g.getTsType(TplType{Name: t.Name, Kind: t.Kind})

	switch {
	case t.Name == "bool" && !t.IsList:
		return fmt.Sprintf("rm.boolField<%s>(%q, %q)", structName, key, label)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "text":
		return fmt.Sprintf("rm.textField<%s>(%q, %q)", structName, key, label)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "bin":
		return fmt.Sprintf("rm.binField<%s>(%q, %q)", structName, key, label)
	case t.Kind == TplKindBase && !t.IsList:
		return fmt.Sprintf("rm.scalarField<%s, %s>(%q, %q, %s, rt.%s, rt.%s, %s)", structName, tsType, key, label, g.primitiveDefault(t.Name), g.primitiveGetter(t.Name), g.primitiveSetter(t.Name), g.primitiveEq(t.Name))
	case t.Kind == TplKindEnum && !t.IsList:
		return fmt.Sprintf("rm.enumField<%s, %s>(%q, %q, Default%s, Is%s, IsAssignable%s, Normalize%s)", structName, typeName, key, label, typeName, typeName, typeName, typeName)
	case t.Kind == TplKindStruct && !t.IsList:
		return fmt.Sprintf("rm.structField<%s, %s>(%q, %q, isZero%s, read%s, set%s, validate%s, eq%s)", structName, typeName, key, label, typeName, typeName, typeName, typeName, typeName)
	case t.IsList && t.Name == "bool":
		return fmt.Sprintf("rm.boolListField<%s>(%q, %q)", structName, key, label)
	case t.IsList && t.Name == "text":
		return fmt.Sprintf("rm.textListField<%s>(%q, %q)", structName, key, label)
	case t.IsList && t.Name == "bin":
		return fmt.Sprintf("rm.binListField<%s>(%q, %q)", structName, key, label)
	case t.IsList && t.Kind == TplKindBase:
		return fmt.Sprintf("rm.zeroListField<%s, %s>(%q, %q, %s, rt.%s, rt.%s, %s)", structName, tsType, key, label, g.primitiveDefault(t.Name), g.primitiveGetter(t.Name), g.primitiveSetter(t.Name), g.primitiveEq(t.Name))
	case t.IsList && t.Kind == TplKindEnum:
		return fmt.Sprintf("rm.defaultListField<%s, %s>(%q, %q, Default%s, get%sListBody, set%sListBody, (item) => IsAssignable%s(item as any) ? undefined : new Error(`非法枚举值: ${item as any}`), (item) => IsDefault%s(item as any), (left, right) => Normalize%s(left as any) === Normalize%s(right as any))", structName, typeName, key, label, typeName, typeName, typeName, typeName, typeName, typeName, typeName)
	default:
		return fmt.Sprintf("rm.defaultListField<%s, %s>(%q, %q, new%s, get%sListBody, set%sListBody, validate%s, isZero%s, eq%s)", structName, typeName, key, label, typeName, typeName, typeName, typeName, typeName, typeName)
	}
}

func (g *TsGenerator) listGetExpr(t TplType, stateVar string) string {
	switch t.Name {
	case "bool":
		return fmt.Sprintf("rt.getBoolList(buf, %s)", stateVar)
	case "text":
		return fmt.Sprintf("rt.getTextList(buf, %s)", stateVar)
	case "bin":
		return fmt.Sprintf("rt.getBinList(buf, %s)", stateVar)
	default:
		if t.Kind == TplKindBase {
			return fmt.Sprintf("rt.getZeroList<%s>(buf, %s, %s, rt.%s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), stateVar, g.primitiveDefault(t.Name), g.primitiveGetter(t.Name))
		}
		return fmt.Sprintf("rt.getDefaultList<%s>(buf, %s, %s, %s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), stateVar, g.bitmapDefaultFactory(t), g.bitmapGetter(t))
	}
}

func (g *TsGenerator) listSetExpr(t TplType, ref, stateVar string) string {
	switch t.Name {
	case "bool":
		return fmt.Sprintf("rt.setBoolList(buf, %s, %s)", stateVar, ref)
	case "text":
		return fmt.Sprintf("rt.setTextList(buf, %s, %s)", stateVar, ref)
	case "bin":
		return fmt.Sprintf("rt.setBinList(buf, %s, %s)", stateVar, ref)
	default:
		if t.Kind == TplKindBase {
			return fmt.Sprintf("rt.setZeroList<%s>(buf, %s, %s, %s, rt.%s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), stateVar, ref, g.primitiveDefault(t.Name), g.primitiveSetter(t.Name))
		}
		return fmt.Sprintf("rt.setDefaultList<%s>(buf, %s, %s, %s, %s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), stateVar, ref, g.bitmapIsDefault(t), g.bitmapSetter(t))
	}
}

func (g *TsGenerator) bitmapDefaultFactory(t TplType) string {
	switch t.Kind {
	case TplKindEnum:
		return fmt.Sprintf("() => _.Default%s()", PascalCase(t.Name))
	case TplKindStruct:
		return fmt.Sprintf("() => _.new%s()", PascalCase(t.Name))
	default:
		return fmt.Sprintf("() => %s", g.primitiveDefault(t.Name))
	}
}

func (g *TsGenerator) bitmapIsDefault(t TplType) string {
	switch t.Kind {
	case TplKindEnum:
		return fmt.Sprintf("(item) => _.IsDefault%s(item as any)", PascalCase(t.Name))
	case TplKindStruct:
		return fmt.Sprintf("(item) => _.isZero%s(item)", PascalCase(t.Name))
	default:
		if t.Name == "i64" || t.Name == "u64" {
			return "(item) => item === 0n"
		}
		return "(item) => item === 0"
	}
}

func (g *TsGenerator) bitmapGetter(t TplType) string {
	switch t.Kind {
	case TplKindEnum:
		return fmt.Sprintf("(buf) => { const [value, err] = rt.getU8(buf); if (err !== null) return [_.Default%s(), rt.errU(err)]; const item = value as _.%s; if (!_.Is%s(item)) return [_.Default%s(), new Error(`非法枚举值: ${item}`)]; return [item, undefined]; }", PascalCase(t.Name), PascalCase(t.Name), PascalCase(t.Name), PascalCase(t.Name))
	case TplKindStruct:
		return fmt.Sprintf("(buf) => _.read%s(buf)", PascalCase(t.Name))
	default:
		return fmt.Sprintf("(buf) => rt.resultU(...rt.%s(buf))", g.primitiveGetter(t.Name))
	}
}

func (g *TsGenerator) bitmapSetter(t TplType) string {
	switch t.Kind {
	case TplKindEnum:
		return fmt.Sprintf("(buf, item) => { if (!_.Is%s(item as any)) return new Error(`非法枚举值: ${item}`); return rt.setU8(buf, item as any); }", PascalCase(t.Name))
	case TplKindStruct:
		return fmt.Sprintf("(buf, item) => _.set%s(buf, item)", PascalCase(t.Name))
	default:
		return fmt.Sprintf("(buf, item) => rt.errU(rt.%s(buf, item as any))", g.primitiveSetter(t.Name))
	}
}

func (g *TsGenerator) eqExpr(t TplType, left, right string) string {
	switch {
	case t.Kind == TplKindStruct && t.IsList:
		return fmt.Sprintf("rt.eqList(%s, %s, _.eq%s)", left, right, PascalCase(t.Name))
	case t.Kind == TplKindStruct:
		return fmt.Sprintf("_.eq%s(%s, %s)", PascalCase(t.Name), left, right)
	case t.Kind == TplKindEnum && t.IsList:
		return fmt.Sprintf("_.eq%sList(%s as any, %s as any)", PascalCase(t.Name), left, right)
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("_.eq%sValue(%s as any, %s as any)", PascalCase(t.Name), left, right)
	case t.Kind == TplKindBase && t.IsList && t.Name == "bin":
		return fmt.Sprintf("rt.eqBinList(%s, %s)", left, right)
	case t.Kind == TplKindBase && t.IsList:
		return fmt.Sprintf("rt.eqList(%s, %s, %s)", left, right, g.primitiveEq(t.Name))
	default:
		return fmt.Sprintf("%s(%s, %s)", g.primitiveEq(t.Name), left, right)
	}
}
