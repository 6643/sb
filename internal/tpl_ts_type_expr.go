package internal

import "fmt"

func (g *TsGenerator) fieldType(t TplType) string {
	base := g.getTsType(TplType{Name: t.Name, Kind: t.Kind})
	switch t.Kind {
	case TplKindEnum, TplKindStruct:
		base = "_." + base
	}
	if t.IsList {
		return base + "[]"
	}
	if t.Kind == TplKindStruct {
		return base + " | null"
	}
	return base
}

func (g *TsGenerator) defaultValue(t TplType) string {
	if t.IsList {
		return "[]"
	}
	switch t.Kind {
	case TplKindBase:
		return g.getTsValue(t.Name)
	case TplKindEnum:
		return fmt.Sprintf("_.Default%s()", PascalCase(t.Name))
	default:
		return "null"
	}
}

func (g *TsGenerator) headerBits(st TplStruct) int {
	total := 0
	for _, field := range st.Fields {
		total += g.tagWidth(field.Type)
	}
	return total
}

func (g *TsGenerator) tagWidth(t TplType) int {
	if t.Name == "bool" && !t.IsList {
		return 1
	}
	if t.Kind == TplKindBase && (t.Name == "text" || t.Name == "bin" || t.IsList) {
		return 2
	}
	if t.IsList {
		return 2
	}
	return 1
}

func (g *TsGenerator) nonDefaultExpr(t TplType, ref string) string {
	switch {
	case t.IsList:
		return fmt.Sprintf("!rt.isArrayValue(%s) || %s.length !== 0", ref, ref)
	case t.Kind == TplKindStruct:
		return fmt.Sprintf("!_.isZero%s(%s as any)", PascalCase(t.Name), ref)
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("!_.IsDefault%s(%s as any)", PascalCase(t.Name), ref)
	case t.Name == "text":
		return fmt.Sprintf("!rt.isStringValue(%s) || %s !== \"\"", ref, ref)
	case t.Name == "bin":
		return fmt.Sprintf("!rt.isBinValue(%s) || %s.byteLength !== 0", ref, ref)
	case t.Name == "bool":
		return fmt.Sprintf("%s !== false", ref)
	case t.Name == "i64", t.Name == "u64":
		return fmt.Sprintf("%s !== 0n", ref)
	default:
		return fmt.Sprintf("%s !== 0", ref)
	}
}

func (g *TsGenerator) primitiveGetter(name string) string {
	switch name {
	case "i8":
		return "getI8"
	case "u8":
		return "getU8"
	case "i16":
		return "getI16"
	case "u16":
		return "getU16"
	case "i32":
		return "getI32"
	case "u32":
		return "getU32"
	case "i64":
		return "getI64"
	case "u64":
		return "getU64"
	case "f32":
		return "getF32"
	case "f64":
		return "getF64"
	default:
		return ""
	}
}

func (g *TsGenerator) primitiveSetter(name string) string {
	switch name {
	case "i8":
		return "setI8"
	case "u8":
		return "setU8"
	case "i16":
		return "setI16"
	case "u16":
		return "setU16"
	case "i32":
		return "setI32"
	case "u32":
		return "setU32"
	case "i64":
		return "setI64"
	case "u64":
		return "setU64"
	case "f32":
		return "setF32"
	case "f64":
		return "setF64"
	default:
		return ""
	}
}

func (g *TsGenerator) primitiveDefault(name string) string {
	switch name {
	case "i64", "u64":
		return "0n"
	default:
		return g.getTsValue(name)
	}
}

func (g *TsGenerator) primitiveEq(name string) string {
	switch name {
	case "i8":
		return "rt.eqI8"
	case "u8":
		return "rt.eqU8"
	case "i16":
		return "rt.eqI16"
	case "u16":
		return "rt.eqU16"
	case "i32":
		return "rt.eqI32"
	case "u32":
		return "rt.eqU32"
	case "i64":
		return "rt.eqI64"
	case "u64":
		return "rt.eqU64"
	case "f32":
		return "rt.eqF32"
	case "f64":
		return "rt.eqF64"
	case "bool":
		return "rt.eqBool"
	case "text":
		return "rt.eqText"
	case "bin":
		return "rt.eqBin"
	default:
		return ""
	}
}
