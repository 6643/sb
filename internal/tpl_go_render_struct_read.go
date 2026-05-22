package internal

func (g *GoGenerator) renderStructReadAndEq(w *sourceWriter, st TplStruct, structName, readFunc, eqFunc, getFunc, defaultFunc string) {
	w.Linef("func %s(buf *bytes.Buffer) (*%s, error) {", readFunc, structName)
	w.Linef("\ts := %s()", defaultFunc)
	w.Linef("\treturn s, %s(buf, s)", getFunc)
	w.Line("}")
	w.Blank()
	w.Linef("func %s(a, b *%s) bool { return EqStruct(%sMeta, a, b) }", eqFunc, structName, CamelCase(st.Name))
	w.Blank()
}
