package internal

import "fmt"

func (g *GoGenerator) directRead(w *sourceWriter, t TplType, target, bufVar, onErr string) {
	w.Line("\t{")
	switch {
	case t.Name == "bool" && !t.IsList:
		w.Linef("\t\tvalue, err := GetU8(%s)", bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\tif value > 1 { %s }", onErr)
		w.Linef("\t\t%s = value != 0", target)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "text":
		w.Linef("\t\tstate, err := GetU8(%s)", bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\tvalue, err := GetText(%s, state)", bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\t%s = value", target)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "bin":
		w.Linef("\t\tstate, err := GetU8(%s)", bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\tvalue, err := GetBinInto(%s, state, nil)", bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\t%s = value", target)
	case t.IsList:
		w.Linef("\t\tstate, err := GetU8(%s)", bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\tvalue, err := %s", g.directListGetExpr(t, "state", bufVar))
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\t%s = value", target)
	case t.Kind == TplKindStruct:
		w.Linef("\t\tvalue, err := %s(%s)", g.structReadName(t.Name), bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\t%s = value", target)
	case t.Kind == TplKindEnum:
		typeName := PascalCase(t.Name)
		w.Linef("\t\tvalue, err := GetU8(%s)", bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\titem := %s(value)", typeName)
		w.Linef("\t\tif !%s(item) { %s }", g.enumIsName(typeName), onErr)
		w.Linef("\t\t%s = item", target)
	default:
		_, getter := g.primitiveGetter(t.Name)
		w.Linef("\t\tvalue, err := %s(%s)", getter, bufVar)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\t%s = value", target)
	}
	w.Line("\t}")
}

func (g *GoGenerator) directWrite(w *sourceWriter, t TplType, ref, bufVar, onErr string) {
	w.Line("\t{")
	switch {
	case t.Name == "bool" && !t.IsList:
		w.Line("\t\tvar value uint8")
		w.Linef("\t\tif %s { value = 1 }", ref)
		w.Linef("\t\tif err := SetU8(%s, value); err != nil { %s }", bufVar, onErr)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "text":
		w.Linef("\t\tstate, err := TextState(len(%s))", ref)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\tif err := SetU8(%s, state); err != nil { %s }", bufVar, onErr)
		w.Linef("\t\tif err := SetText(%s, state, %s); err != nil { %s }", bufVar, ref, onErr)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "bin":
		w.Linef("\t\tstate, err := BinState(len(%s))", ref)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\tif err := SetU8(%s, state); err != nil { %s }", bufVar, onErr)
		w.Linef("\t\tif err := SetBin(%s, state, %s); err != nil { %s }", bufVar, ref, onErr)
	case t.IsList:
		w.Linef("\t\tstate, err := ListCountState(len(%s))", ref)
		w.Linef("\t\tif err != nil { %s }", onErr)
		w.Linef("\t\tif err := SetU8(%s, state); err != nil { %s }", bufVar, onErr)
		w.Linef("\t\tif err := %s; err != nil { %s }", g.directListSetExpr(t, ref, "state", bufVar), onErr)
	case t.Kind == TplKindStruct:
		w.Linef("\t\tif err := %s(%s, %s); err != nil { %s }", g.structSetName(t.Name), bufVar, ref, onErr)
	case t.Kind == TplKindEnum:
		typeName := PascalCase(t.Name)
		w.Linef("\t\tif !%s(%s) { %s }", g.enumIsAssignableName(typeName), ref, onErr)
		w.Linef("\t\tif err := SetU8(%s, uint8(%s(%s))); err != nil { %s }", bufVar, g.enumNormalizeName(typeName), ref, onErr)
	default:
		w.Linef("\t\tif err := %s(%s, %s); err != nil { %s }", g.primitiveSetter(t.Name), bufVar, ref, onErr)
	}
	w.Line("\t}")
}

func (g *GoGenerator) directListGetExpr(t TplType, stateVar, bufVar string) string {
	switch {
	case t.Name == "bool":
		return fmt.Sprintf("GetBoolListInto(%s, %s, nil)", bufVar, stateVar)
	case t.Name == "text":
		return fmt.Sprintf("GetTextListInto(%s, %s, nil)", bufVar, stateVar)
	case t.Name == "bin":
		return fmt.Sprintf("GetBinListInto(%s, %s, nil)", bufVar, stateVar)
	case t.Kind == TplKindBase:
		_, getter := g.primitiveGetter(t.Name)
		return fmt.Sprintf("GetZeroListInto(%s, %s, nil, %s)", bufVar, stateVar, getter)
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("get%sListBody(%s, %s)", PascalCase(t.Name), bufVar, stateVar)
	default:
		return fmt.Sprintf("get%sListBody(%s, %s)", PascalCase(t.Name), bufVar, stateVar)
	}
}

func (g *GoGenerator) directListSetExpr(t TplType, ref, stateVar, bufVar string) string {
	switch {
	case t.Name == "bool":
		return fmt.Sprintf("SetBoolList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Name == "text":
		return fmt.Sprintf("SetTextList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Name == "bin":
		return fmt.Sprintf("SetBinList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Kind == TplKindBase:
		width, _ := goBaseEncodedWidth(t.Name)
		return fmt.Sprintf("SetZeroList(%s, %s, %s, %d, %s)", bufVar, stateVar, ref, width, g.primitiveSetter(t.Name))
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("set%sListBody(%s, %s, %s)", PascalCase(t.Name), bufVar, stateVar, ref)
	default:
		return fmt.Sprintf("set%sListBody(%s, %s, %s)", PascalCase(t.Name), bufVar, stateVar, ref)
	}
}
