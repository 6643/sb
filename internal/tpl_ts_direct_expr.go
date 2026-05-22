package internal

import "fmt"

func (g *TsGenerator) tsDirectRead(w *sourceWriter, t TplType, target, bufVar, errReturn string) {
	w.Line("        {")
	switch {
	case t.Name == "bool" && !t.IsList:
		w.Linef("            const [value, err] = rt.getU8(%s);", bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            if (value > 1) return %s;", errReturn)
		w.Linef("            %s = (value !== 0) as any;", target)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "text":
		w.Linef("            const [state, errState] = rt.getU8(%s);", bufVar)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const [value, err] = rt.getText(%s, state);", bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            %s = value as any;", target)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "bin":
		w.Linef("            const [state, errState] = rt.getU8(%s);", bufVar)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const [value, err] = rt.getBin(%s, state);", bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            %s = value as any;", target)
	case t.IsList:
		w.Linef("            const [state, errState] = rt.getU8(%s);", bufVar)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const [value, err] = %s;", g.tsDirectListGetExpr(t, "state", bufVar))
		if t.Name == "bool" || t.Name == "text" || t.Name == "bin" {
			w.Linef("            if (err !== null) return %s;", errReturn)
		} else {
			w.Linef("            if (err !== undefined) return %s;", errReturn)
		}
		w.Linef("            %s = value as any;", target)
	case t.Kind == TplKindStruct:
		w.Linef("            const [value, err] = read%s(%s);", PascalCase(t.Name), bufVar)
		w.Linef("            if (err !== undefined) return %s;", errReturn)
		w.Linef("            %s = value as any;", target)
	case t.Kind == TplKindEnum:
		typeName := PascalCase(t.Name)
		w.Linef("            const [value, err] = rt.getU8(%s);", bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            const item = value as %s;", typeName)
		w.Linef("            if (!Is%s(item)) return %s;", typeName, errReturn)
		w.Linef("            %s = item as any;", target)
	default:
		w.Linef("            const [value, err] = rt.%s(%s);", g.primitiveGetter(t.Name), bufVar)
		w.Linef("            if (err !== null) return %s;", errReturn)
		w.Linef("            %s = value as any;", target)
	}
	w.Line("        }")
}

func (g *TsGenerator) tsDirectWrite(w *sourceWriter, t TplType, ref, bufVar, errReturn string) {
	w.Line("        {")
	switch {
	case t.Name == "bool" && !t.IsList:
		w.Linef("            const encodeErr = rt.setU8(%s, %s ? 1 : 0);", bufVar, ref)
		w.Linef("            if (encodeErr !== null) return %s;", errReturn)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "text":
		w.Linef("            const [state, errState] = rt.textState(%s);", ref)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const err0 = rt.setU8(%s, state);", bufVar)
		w.Linef("            if (err0 !== null) return %s;", errReturn)
		w.Linef("            const err1 = rt.setText(%s, state, %s);", bufVar, ref)
		w.Linef("            if (err1 !== null) return %s;", errReturn)
	case t.Kind == TplKindBase && !t.IsList && t.Name == "bin":
		w.Linef("            const [state, errState] = rt.binState(%s.byteLength);", ref)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const err0 = rt.setU8(%s, state);", bufVar)
		w.Linef("            if (err0 !== null) return %s;", errReturn)
		w.Linef("            const err1 = rt.setBin(%s, state, %s);", bufVar, ref)
		w.Linef("            if (err1 !== null) return %s;", errReturn)
	case t.IsList:
		w.Linef("            const [state, errState] = rt.listCountState(%s.length);", ref)
		w.Linef("            if (errState !== null) return %s;", errReturn)
		w.Linef("            const err0 = rt.setU8(%s, state);", bufVar)
		w.Linef("            if (err0 !== null) return %s;", errReturn)
		w.Linef("            const err1 = %s;", g.tsDirectListSetExpr(t, ref, "state", bufVar))
		if t.Name == "bool" || t.Name == "text" || t.Name == "bin" {
			w.Linef("            if (err1 !== null) return %s;", errReturn)
		} else {
			w.Linef("            if (err1 !== undefined) return %s;", errReturn)
		}
	case t.Kind == TplKindStruct:
		w.Linef("            const encodeErr = set%s(%s, %s);", PascalCase(t.Name), bufVar, ref)
		w.Linef("            if (encodeErr !== undefined) return %s;", errReturn)
	case t.Kind == TplKindEnum:
		typeName := PascalCase(t.Name)
		w.Linef("            if (!IsAssignable%s(%s as any)) return %s;", typeName, ref, errReturn)
		w.Linef("            const encodeErr = rt.setU8(%s, Normalize%s(%s as any) as any);", bufVar, typeName, ref)
		w.Linef("            if (encodeErr !== null) return %s;", errReturn)
	default:
		w.Linef("            const encodeErr = rt.%s(%s, %s as any);", g.primitiveSetter(t.Name), bufVar, ref)
		w.Linef("            if (encodeErr !== null) return %s;", errReturn)
	}
	w.Line("        }")
}

func (g *TsGenerator) tsDirectListGetExpr(t TplType, stateVar, bufVar string) string {
	switch {
	case t.Name == "bool":
		return fmt.Sprintf("rt.getBoolList(%s, %s)", bufVar, stateVar)
	case t.Name == "text":
		return fmt.Sprintf("rt.getTextList(%s, %s)", bufVar, stateVar)
	case t.Name == "bin":
		return fmt.Sprintf("rt.getBinList(%s, %s)", bufVar, stateVar)
	case t.Kind == TplKindBase:
		return fmt.Sprintf("rt.getDefaultList<%s>(%s, %s, %s, %s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), bufVar, stateVar, g.bitmapDefaultFactory(t), g.bitmapGetter(t))
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("get%sListBody(%s, %s)", PascalCase(t.Name), bufVar, stateVar)
	default:
		return fmt.Sprintf("get%sListBody(%s, %s)", PascalCase(t.Name), bufVar, stateVar)
	}
}

func (g *TsGenerator) tsDirectListSetExpr(t TplType, ref, stateVar, bufVar string) string {
	switch {
	case t.Name == "bool":
		return fmt.Sprintf("rt.setBoolList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Name == "text":
		return fmt.Sprintf("rt.setTextList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Name == "bin":
		return fmt.Sprintf("rt.setBinList(%s, %s, %s)", bufVar, stateVar, ref)
	case t.Kind == TplKindBase:
		return fmt.Sprintf("rt.setDefaultList<%s>(%s, %s, %s, %s, %s)", g.getTsType(TplType{Name: t.Name, Kind: t.Kind}), bufVar, stateVar, ref, g.bitmapIsDefault(t), g.bitmapSetter(t))
	case t.Kind == TplKindEnum:
		return fmt.Sprintf("set%sListBody(%s, %s, %s)", PascalCase(t.Name), bufVar, stateVar, ref)
	default:
		return fmt.Sprintf("set%sListBody(%s, %s, %s)", PascalCase(t.Name), bufVar, stateVar, ref)
	}
}

func (g *TsGenerator) readErrSuffix(t TplType) string {
	if t.Name == "bool" && !t.IsList {
		return "State"
	}
	if g.tagWidth(t) == 1 {
		return "Present"
	}
	return "State"
}
