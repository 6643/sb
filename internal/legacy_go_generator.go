package internal

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Generate 生成 Go 代码。
func GenerateLegacyGo(schema *IRSchema, cfg Config) error {
	dir, err := newGeneratedDir(filepath.Join(cfg.GoDir, "sb"))
	if err != nil {
		return err
	}

	if err := legacyGoGenerateTypeFile(dir); err != nil {
		return err
	}
	if err := legacyGoGenerateEnumFile(dir, schema.Enums); err != nil {
		return err
	}
	if err := legacyGoGenerateStructFiles(dir, schema.Structs, cfg.GoTag); err != nil {
		return err
	}
	if err := legacyGoGenerateAPIStubFiles(dir, schema.APIs); err != nil {
		return err
	}
	if err := legacyGoGenerateRPCFile(dir, schema.APIs); err != nil {
		return err
	}
	return nil
}

func legacyGoGenerateTypeFile(dir *generatedDir) error {
	content := `package sb

type RpcErrCode int

const (
	RpcOk       RpcErrCode = 200
	RpcNoConn   RpcErrCode = 0
	RpcTimeout  RpcErrCode = 408
	RpcReqErr   RpcErrCode = 400
	RpcRespErr  RpcErrCode = 500
	RpcNotAuth  RpcErrCode = 401
	RpcNotExist RpcErrCode = 404
)
`
	return dir.Write("type.go", []byte(content), 0644)
}

func legacyGoGenerateEnumFile(dir *generatedDir, enums []IREnum) error {
	var b bytes.Buffer
	b.WriteString("package sb\n\n")
	if len(enums) == 0 {
		b.WriteString("// No enums defined.\n")
		return dir.Write("enum.go", b.Bytes(), 0644)
	}

	for _, e := range enums {
		enumName := PascalCase(e.Name)
		if e.Note != "" {
			b.WriteString(RenderLineCommentWithHead("// ", enumName+" ", e.Note))
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "type %s uint8\n\n", enumName)
		b.WriteString("const (\n")
		for _, m := range e.Members {
			memberName := enumName + PascalCase(m.Name)
			if m.Note != "" {
				b.WriteString(RenderLineComment("\t// ", m.Note))
				b.WriteByte('\n')
			}
			fmt.Fprintf(&b, "\t%s %s = %d\n", memberName, enumName, m.ID)
		}
		b.WriteString(")\n\n")
	}

	return dir.Write("enum.go", b.Bytes(), 0644)
}

func legacyGoGenerateStructFiles(dir *generatedDir, structs []IRStruct, goTag string) error {
	for _, st := range structs {
		var b bytes.Buffer
		fmt.Fprintln(&b, "package sb")
		fmt.Fprintln(&b)
		if st.Note != "" {
			b.WriteString(RenderLineCommentWithHead("// ", PascalCase(st.Name)+" ", st.Note))
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "type %s struct {\n", PascalCase(st.Name))
		for _, f := range st.Fields {
			fieldName := PascalCase(f.Name)
			fieldType := legacyGoType(f.Type)
			tag := legacyGoTagString(goTag, f)
			if f.Note != "" {
				b.WriteString(RenderLineComment("\t// ", f.Note))
				b.WriteByte('\n')
			}
			if tag == "" {
				fmt.Fprintf(&b, "\t%s %s\n", fieldName, fieldType)
				continue
			}
			fmt.Fprintf(&b, "\t%s %s %s\n", fieldName, fieldType, tag)
		}
		fmt.Fprintln(&b, "}")

		filename := "struct_" + SnakeCase(st.Name) + ".go"
		if err := dir.Write(filename, b.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

func legacyGoGenerateAPIStubFiles(dir *generatedDir, apis []IRAPI) error {
	for _, api := range apis {
		filename := "api." + api.Name + ".go"
		var b bytes.Buffer
		fmt.Fprintln(&b, "package sb")
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "import (\n\t\"context\"\n)")
		fmt.Fprintln(&b)

		funcName := SnakeCase(api.Name)
		args := legacyGoLogicArgs(api.Args)
		if api.Result.Kind == IRKindNil {
			fmt.Fprintf(&b, "func %s(ctx context.Context%s) (errCode RpcErrCode) {\n", funcName, args)
			fmt.Fprintln(&b, "\treturn RpcRespErr")
			fmt.Fprintln(&b, "}")
		} else {
			retType := legacyGoType(api.Result)
			fmt.Fprintf(&b, "func %s(ctx context.Context%s) (result %s, errCode RpcErrCode) {\n", funcName, args, retType)
			fmt.Fprintf(&b, "\treturn %s, RpcRespErr\n", legacyGoZeroValue(api.Result))
			fmt.Fprintln(&b, "}")
		}

		if err := dir.WriteIfAbsent(filename, b.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

func legacyGoGenerateRPCFile(dir *generatedDir, apis []IRAPI) error {
	var b bytes.Buffer
	fmt.Fprintln(&b, "package sb")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "import (\n\t\"context\"\n)")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "type Client struct{}")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "func NewClient() *Client { return &Client{} }")
	fmt.Fprintln(&b)

	groups := map[string][]string{}
	for _, api := range apis {
		group := APIGroup(api.Name)
		groups[group] = append(groups[group], api.Name)
	}
	for group := range groups {
		sort.Strings(groups[group])
	}

	for _, api := range apis {
		methodName := PascalCase(api.Name)
		args := legacyGoLogicArgs(api.Args)
		if api.Result.Kind == IRKindNil {
			fmt.Fprintf(&b, "func Call%s(c *Client, ctx context.Context%s) (errCode RpcErrCode) {\n", methodName, args)
			fmt.Fprintln(&b, "\treturn RpcRespErr")
			fmt.Fprintln(&b, "}")
			fmt.Fprintln(&b)
			continue
		}

		retType := legacyGoType(api.Result)
		fmt.Fprintf(&b, "func Call%s(c *Client, ctx context.Context%s) (result %s, errCode RpcErrCode) {\n", methodName, args, retType)
		fmt.Fprintf(&b, "\treturn %s, RpcRespErr\n", legacyGoZeroValue(api.Result))
		fmt.Fprintln(&b, "}")
		fmt.Fprintln(&b)
	}

	if len(groups) > 0 {
		keys := make([]string, 0, len(groups))
		for k := range groups {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fmt.Fprintln(&b, "// API 路径分组常量")
		for _, k := range keys {
			fmt.Fprintf(&b, "var API%s = []string{\n", PascalCase(k))
			for _, name := range groups[k] {
				fmt.Fprintf(&b, "\t\"%s\",\n", name)
			}
			fmt.Fprintln(&b, "}")
		}
	}

	return dir.Write("rpc.go", b.Bytes(), 0644)
}

func legacyGoLogicArgs(args []IRAPIArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, a := range args {
		parts = append(parts, fmt.Sprintf(", %s %s", CamelCase(a.Name), legacyGoType(a.Type)))
	}
	return strings.Join(parts, "")
}

func legacyGoType(t IRType) string {
	prefix := ""
	if t.IsList {
		prefix = "[]"
	}

	switch t.Kind {
	case IRKindBase:
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
	case IRKindEnum:
		return prefix + PascalCase(t.Name)
	case IRKindStruct:
		return prefix + "*" + PascalCase(t.Name)
	}

	return prefix + PascalCase(t.Name)
}

func legacyGoZeroValue(t IRType) string {
	if t.IsList {
		return "nil"
	}
	if t.Kind == IRKindStruct || t.Kind == IRKindNil {
		return "nil"
	}
	if t.Kind == IRKindEnum {
		return "0"
	}

	switch t.Name {
	case "text":
		return "\"\""
	case "bin":
		return "nil"
	case "bool":
		return "false"
	case "f32", "f64":
		return "0"
	default:
		return "0"
	}
}

func legacyGoTagString(goTag string, f IRField) string {
	if strings.TrimSpace(goTag) == "" {
		return ""
	}

	keys := strings.Split(goTag, ",")
	val := f.Tag
	if strings.TrimSpace(val) == "" {
		val = SnakeCase(f.Name)
	}

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s:\"%s\"", k, val))
	}
	if len(parts) == 0 {
		return ""
	}
	return "`" + strings.Join(parts, " ") + "`"
}
