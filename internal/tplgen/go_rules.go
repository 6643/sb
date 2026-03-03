package tplgen

import (
	"fmt"
	"sb/internal/util"
	"strings"
)

func (g *GoGenerator) getGoLogicType(t Type) string {
	if t.Name == "nil" {
		return ""
	}
	return g.getGoType(t)
}

func (g *GoGenerator) getGoRpcType(t Type) string {
	if t.Name == "nil" {
		return ""
	}
	if t.IsList {
		return util.PascalCase(t.Name) + "List"
	}
	if t.Kind == KindBase {
		return util.PascalCase(t.Name)
	}
	if t.Kind == KindEnum {
		return "U8"
	}
	return util.PascalCase(t.Name)
}

func (g *GoGenerator) getGoType(t Type) string {
	if t.Name == "nil" {
		return ""
	}

	prefix := ""
	if t.IsList {
		prefix = "[]"
	}

	switch t.Name {
	case "i8":
		return prefix + "int8"
	case "u8":
		return prefix + "uint8"
	case "i16":
		return prefix + "int16"
	case "u16":
		return prefix + "uint16"
	case "i32":
		return prefix + "int32"
	case "u32":
		return prefix + "uint32"
	case "i64":
		return prefix + "int64"
	case "u64":
		return prefix + "uint64"
	case "f32":
		return prefix + "float32"
	case "f64":
		return prefix + "float64"
	case "bool":
		return prefix + "bool"
	case "text":
		return prefix + "string"
	case "bin":
		return prefix + "[]byte"
	}

	if t.Kind == KindStruct {
		return prefix + "*" + util.PascalCase(t.Name)
	}
	return prefix + util.PascalCase(t.Name)
}

func (g *GoGenerator) getGoValue(name string) string {
	switch name {
	case "text":
		return "\"\""
	case "bin", "nil":
		return "nil"
	case "bool":
		return "false"
	case "f32", "f64":
		return "0.0"
	default:
		return "0"
	}
}

func (g *GoGenerator) getGoTag(field StructField) string {
	if g.Config.GoTag == "" {
		return ""
	}

	val := field.Tag
	if val == "" {
		val = util.SnakeCase(field.Name)
	}

	parts := strings.Split(g.Config.GoTag, ",")
	list := make([]string, 0, len(parts))
	for _, k := range parts {
		key := strings.TrimSpace(k)
		if key == "" {
			continue
		}
		list = append(list, fmt.Sprintf("%s:\"%s\"", key, val))
	}
	if len(list) == 0 {
		return ""
	}
	return "`" + strings.Join(list, " ") + "`"
}
