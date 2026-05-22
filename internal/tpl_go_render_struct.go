package internal

func (g *GoGenerator) renderStruct(w *sourceWriter, st TplStruct) {
	structName := PascalCase(st.Name)
	resetFunc := g.structResetName(structName)
	defaultFunc := g.structDefaultName(structName)
	isZeroFunc := g.structIsZeroName(structName)
	validateFunc := g.structValidateName(structName)
	sizeFunc := g.structSizeName(st.Name)
	getFunc := g.structGetName(st.Name)
	setFunc := g.structSetName(st.Name)
	readFunc := g.structReadName(st.Name)
	eqFunc := g.structEqName(st.Name)
	sizeValidatedFunc := g.structSizeValidatedName(st.Name)
	setValidatedFunc := g.structSetValidatedName(st.Name)
	headerBits := g.headerBits(st)
	headerSize := (headerBits + 7) / 8
	headerWidthsName := g.structHeaderWidthsName(structName)
	headerFieldCount := len(st.Fields)

	g.renderStructTypeAndMeta(w, st, structName, resetFunc, defaultFunc, isZeroFunc, validateFunc, headerWidthsName)
	g.renderStructSize(w, st, structName, sizeFunc, sizeValidatedFunc, validateFunc, headerBits)
	g.renderStructGet(w, st, structName, getFunc, resetFunc, headerBits, headerWidthsName, headerFieldCount)
	g.renderStructSet(w, st, structName, setFunc, setValidatedFunc, validateFunc, sizeValidatedFunc, headerSize, headerWidthsName, headerFieldCount)
	g.renderStructReadAndEq(w, st, structName, readFunc, eqFunc, getFunc, defaultFunc)
	g.renderStructListBodyHelpers(w, st)
}
