package internal

import (
	"fmt"
	"strings"
)

func renderDocArgs(args []TplApiArg) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		typeName := arg.Type.Name
		if arg.Type.IsList {
			typeName = "[" + typeName + "]"
		}
		parts = append(parts, fmt.Sprintf("%s %s<br>", arg.Name, typeName))
	}
	return strings.Join(parts, "")
}

func renderDocReturn(result TplType) string {
	if result.Name == "nil" {
		return "Void"
	}
	if result.IsList {
		return "[" + result.Name + "]"
	}
	return result.Name
}

func renderGoExampleArgs(args []TplApiArg) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, renderGoExampleValue(arg.Type))
	}
	if len(parts) == 0 {
		return ""
	}
	return ", " + strings.Join(parts, ", ")
}

func renderGoExampleValue(t TplType) string {
	if t.IsList || t.Kind == TplKindStruct {
		return "nil"
	}
	switch t.Name {
	case "text":
		return "\"\""
	case "bin":
		return "nil"
	case "bool":
		return "false"
	case "i64", "u64":
		return "0"
	default:
		return "0"
	}
}

func renderTsExampleArgs(args []TplApiArg) string {
	parts := make([]string, 0, len(args))
	for _, arg := range args {
		parts = append(parts, renderTsExampleValue(arg.Type))
	}
	return strings.Join(parts, ", ")
}

func renderTsExampleValue(t TplType) string {
	if t.IsList {
		return "[]"
	}
	if t.Kind == TplKindStruct {
		return fmt.Sprintf("sb.new%s()", PascalCase(t.Name))
	}
	switch t.Name {
	case "text":
		return "\"\""
	case "bin":
		return "new Uint8Array(0)"
	case "bool":
		return "false"
	case "i64", "u64":
		return "0n"
	default:
		return "0"
	}
}
