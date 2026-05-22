package internal

import "fmt"

func (g *GoGenerator) headerBits(st TplStruct) int {
	total := 0
	for _, field := range st.Fields {
		total += g.tagWidth(field.Type)
	}
	return total
}

func (g *GoGenerator) tagWidth(t TplType) int {
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

func (g *GoGenerator) nonDefaultExpr(t TplType, ref string) string {
	switch {
	case t.IsList:
		return fmt.Sprintf("len(%s) != 0", ref)
	case t.Kind == TplKindStruct:
		return fmt.Sprintf("!%s(%s)", g.structIsZeroName(PascalCase(t.Name)), ref)
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("!%s(%s)", g.enumIsDefaultName(PascalCase(t.Name)), ref)
	case t.Name == "text":
		return fmt.Sprintf("%s != \"\"", ref)
	case t.Name == "bin":
		return fmt.Sprintf("len(%s) != 0", ref)
	case t.Name == "bool":
		return ref
	default:
		return fmt.Sprintf("%s != 0", ref)
	}
}

func (g *GoGenerator) stateInitLine(field TplStructField, onErr string) string {
	fieldName := PascalCase(field.Name)
	ref := "s." + fieldName
	switch {
	case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "text":
		return fmt.Sprintf("%sState, err := TextState(len(%s)); if err != nil { %s }", CamelCase(fieldName), ref, onErr)
	case field.Type.Kind == TplKindBase && !field.Type.IsList && field.Type.Name == "bin":
		return fmt.Sprintf("%sState, err := BinState(len(%s)); if err != nil { %s }", CamelCase(fieldName), ref, onErr)
	default:
		return fmt.Sprintf("%sState, err := ListCountState(len(%s)); if err != nil { %s }", CamelCase(fieldName), ref, onErr)
	}
}

func (g *GoGenerator) primitiveGetter(name string) (int, string) {
	switch name {
	case "i8":
		return 1, "GetI8"
	case "u8":
		return 1, "GetU8"
	case "i16":
		return 2, "GetI16"
	case "u16":
		return 2, "GetU16"
	case "i32":
		return 4, "GetI32"
	case "u32":
		return 4, "GetU32"
	case "i64":
		return 8, "GetI64"
	case "u64":
		return 8, "GetU64"
	case "f32":
		return 4, "GetF32"
	case "f64":
		return 8, "GetF64"
	default:
		return 0, ""
	}
}

func (g *GoGenerator) primitiveSetter(name string) string {
	switch name {
	case "i8":
		return "SetI8"
	case "u8":
		return "SetU8"
	case "i16":
		return "SetI16"
	case "u16":
		return "SetU16"
	case "i32":
		return "SetI32"
	case "u32":
		return "SetU32"
	case "i64":
		return "SetI64"
	case "u64":
		return "SetU64"
	case "f32":
		return "SetF32"
	case "f64":
		return "SetF64"
	default:
		return ""
	}
}
