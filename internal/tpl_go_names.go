package internal

import (
	"fmt"
	"strings"
)

func (g *GoGenerator) enumDefaultName(typeName string) string { return "default" + typeName }
func (g *GoGenerator) enumIsName(typeName string) string      { return "is" + typeName }
func (g *GoGenerator) enumNormalizeName(typeName string) string {
	return "normalize" + typeName
}
func (g *GoGenerator) enumIsDefaultName(typeName string) string { return "isDefault" + typeName }
func (g *GoGenerator) enumIsAssignableName(typeName string) string {
	return "isAssignable" + typeName
}
func (g *GoGenerator) enumEqValueName(typeName string) string    { return "eq" + typeName + "Value" }
func (g *GoGenerator) enumEqListName(typeName string) string     { return "eq" + typeName + "List" }
func (g *GoGenerator) structResetName(typeName string) string    { return "reset" + typeName }
func (g *GoGenerator) structDefaultName(typeName string) string  { return "default" + typeName }
func (g *GoGenerator) structIsZeroName(typeName string) string   { return "isZero" + typeName }
func (g *GoGenerator) structValidateName(typeName string) string { return "validate" + typeName }
func (g *GoGenerator) structSizeName(rawName string) string {
	suffix := PascalCase(rawName)
	if g.isAPIWrapper(rawName) {
		return "size" + suffix
	}
	return "Size" + suffix
}
func (g *GoGenerator) structGetName(rawName string) string {
	suffix := PascalCase(rawName)
	if g.isAPIWrapper(rawName) {
		return "get" + suffix
	}
	return "Get" + suffix
}
func (g *GoGenerator) structSetName(rawName string) string {
	suffix := PascalCase(rawName)
	if g.isAPIWrapper(rawName) {
		return "set" + suffix
	}
	return "Set" + suffix
}
func (g *GoGenerator) structReadName(rawName string) string {
	suffix := PascalCase(rawName)
	if g.isAPIWrapper(rawName) {
		return "read" + suffix
	}
	return "Read" + suffix
}
func (g *GoGenerator) structEqName(rawName string) string {
	suffix := PascalCase(rawName)
	if g.isAPIWrapper(rawName) {
		return "eq" + suffix
	}
	return "Eq" + suffix
}
func (g *GoGenerator) structSizeValidatedName(rawName string) string {
	return "size" + PascalCase(rawName) + "Validated"
}
func (g *GoGenerator) structSetValidatedName(rawName string) string {
	return "set" + PascalCase(rawName) + "Validated"
}
func (g *GoGenerator) structHeaderWidthsName(typeName string) string {
	return CamelCase(typeName) + "HeaderWidths"
}
func (g *GoGenerator) structHeaderWidthsLiteral(st TplStruct) string {
	widths := make([]string, 0, len(st.Fields))
	for _, field := range st.Fields {
		widths = append(widths, fmt.Sprintf("%d", g.tagWidth(field.Type)))
	}
	return strings.Join(widths, ", ")
}
func (g *GoGenerator) isAPIWrapper(rawName string) bool {
	return len(rawName) > 4 && rawName[:4] == "api_"
}
