package internal

func (g *GoGenerator) renderSharedHelpers(w *sourceWriter) {
	w.Line("func eqBin(a, b []byte) bool { return bytes.Equal(a, b) }")
	w.Line("func eqBinList(a, b [][]byte) bool { return slices.EqualFunc(a, b, bytes.Equal) }")
	w.Line("func eqF32(a, b float32) bool {")
	w.Line("\tif a == b { return true }")
	w.Line("\tif math.IsNaN(float64(a)) && math.IsNaN(float64(b)) { return true }")
	w.Line("\treturn math.Abs(float64(a-b)) < 1e-6")
	w.Line("}")
	w.Line("func eqF64(a, b float64) bool {")
	w.Line("\tif a == b { return true }")
	w.Line("\tif math.IsNaN(a) && math.IsNaN(b) { return true }")
	w.Line("\treturn math.Abs(a-b) < 1e-9")
	w.Line("}")
}
