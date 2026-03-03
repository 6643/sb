package gogen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sb/internal/gen"
	"sb/internal/ir"
	"sb/internal/util"
)

// Generate 生成 Go 代码。
func Generate(schema *ir.Schema, cfg gen.Config) error {
	targetDir := filepath.Join(cfg.GoDir, "sb")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create go dir %s: %w", targetDir, err)
	}

	if err := generateTypeFile(targetDir); err != nil {
		return err
	}
	if err := generateEnumFile(targetDir, schema.Enums); err != nil {
		return err
	}
	if err := generateStructFiles(targetDir, schema.Structs, cfg.GoTag); err != nil {
		return err
	}
	if err := generateAPIStubFiles(targetDir, schema.APIs); err != nil {
		return err
	}
	if err := generateRPCFile(targetDir, schema.APIs); err != nil {
		return err
	}
	return nil
}

func generateTypeFile(targetDir string) error {
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
	return os.WriteFile(filepath.Join(targetDir, "type.go"), []byte(content), 0644)
}

func generateEnumFile(targetDir string, enums []ir.Enum) error {
	var b bytes.Buffer
	b.WriteString("package sb\n\n")
	if len(enums) == 0 {
		b.WriteString("// No enums defined.\n")
		return os.WriteFile(filepath.Join(targetDir, "enum.go"), b.Bytes(), 0644)
	}

	for _, e := range enums {
		enumName := util.PascalCase(e.Name)
		if e.Note != "" {
			fmt.Fprintf(&b, "// %s %s\n", enumName, e.Note)
		}
		fmt.Fprintf(&b, "type %s uint8\n\n", enumName)
		b.WriteString("const (\n")
		for _, m := range e.Members {
			memberName := enumName + util.PascalCase(m.Name)
			if m.Note == "" {
				fmt.Fprintf(&b, "\t%s %s = %d\n", memberName, enumName, m.ID)
				continue
			}
			fmt.Fprintf(&b, "\t%s %s = %d // %s\n", memberName, enumName, m.ID, m.Note)
		}
		b.WriteString(")\n\n")
	}

	return os.WriteFile(filepath.Join(targetDir, "enum.go"), b.Bytes(), 0644)
}

func generateStructFiles(targetDir string, structs []ir.Struct, goTag string) error {
	for _, st := range structs {
		var b bytes.Buffer
		fmt.Fprintln(&b, "package sb")
		fmt.Fprintln(&b)
		if st.Note != "" {
			fmt.Fprintf(&b, "// %s %s\n", util.PascalCase(st.Name), st.Note)
		}
		fmt.Fprintf(&b, "type %s struct {\n", util.PascalCase(st.Name))
		for _, f := range st.Fields {
			fieldName := util.PascalCase(f.Name)
			fieldType := goType(f.Type)
			tag := goTagString(goTag, f)
			if f.Note == "" && tag == "" {
				fmt.Fprintf(&b, "\t%s %s\n", fieldName, fieldType)
				continue
			}
			if f.Note == "" {
				fmt.Fprintf(&b, "\t%s %s %s\n", fieldName, fieldType, tag)
				continue
			}
			if tag == "" {
				fmt.Fprintf(&b, "\t%s %s // %s\n", fieldName, fieldType, f.Note)
				continue
			}
			fmt.Fprintf(&b, "\t%s %s %s // %s\n", fieldName, fieldType, tag, f.Note)
		}
		fmt.Fprintln(&b, "}")

		filename := "struct_" + util.SnakeCase(st.Name) + ".go"
		if err := os.WriteFile(filepath.Join(targetDir, filename), b.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

func generateAPIStubFiles(targetDir string, apis []ir.API) error {
	for _, api := range apis {
		filename := "api." + api.Name + ".go"
		path := filepath.Join(targetDir, filename)
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return err
		}

		var b bytes.Buffer
		fmt.Fprintln(&b, "package sb")
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "import (\n\t\"context\"\n)")
		fmt.Fprintln(&b)

		funcName := util.SnakeCase(api.Name)
		args := goLogicArgs(api.Args)
		if api.Result.Kind == ir.KindNil {
			fmt.Fprintf(&b, "func %s(ctx context.Context%s) (errCode RpcErrCode) {\n", funcName, args)
			fmt.Fprintln(&b, "\treturn RpcRespErr")
			fmt.Fprintln(&b, "}")
		} else {
			retType := goType(api.Result)
			fmt.Fprintf(&b, "func %s(ctx context.Context%s) (result %s, errCode RpcErrCode) {\n", funcName, args, retType)
			fmt.Fprintf(&b, "\treturn %s, RpcRespErr\n", goZeroValue(api.Result))
			fmt.Fprintln(&b, "}")
		}

		if err := os.WriteFile(path, b.Bytes(), 0644); err != nil {
			return err
		}
	}
	return nil
}

func generateRPCFile(targetDir string, apis []ir.API) error {
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
		group := util.APIGroup(api.Name)
		groups[group] = append(groups[group], api.Name)
	}
	for group := range groups {
		sort.Strings(groups[group])
	}

	for _, api := range apis {
		methodName := util.PascalCase(api.Name)
		args := goLogicArgs(api.Args)
		if api.Result.Kind == ir.KindNil {
			fmt.Fprintf(&b, "func (c *Client) %s(ctx context.Context%s) (errCode RpcErrCode) {\n", methodName, args)
			fmt.Fprintln(&b, "\treturn RpcRespErr")
			fmt.Fprintln(&b, "}")
			fmt.Fprintln(&b)
			continue
		}

		retType := goType(api.Result)
		fmt.Fprintf(&b, "func (c *Client) %s(ctx context.Context%s) (result %s, errCode RpcErrCode) {\n", methodName, args, retType)
		fmt.Fprintf(&b, "\treturn %s, RpcRespErr\n", goZeroValue(api.Result))
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
			fmt.Fprintf(&b, "var API%s = []string{\n", util.PascalCase(k))
			for _, name := range groups[k] {
				fmt.Fprintf(&b, "\t\"%s\",\n", name)
			}
			fmt.Fprintln(&b, "}")
		}
	}

	return os.WriteFile(filepath.Join(targetDir, "rpc.go"), b.Bytes(), 0644)
}

func goLogicArgs(args []ir.APIArg) string {
	if len(args) == 0 {
		return ""
	}
	parts := make([]string, 0, len(args))
	for _, a := range args {
		parts = append(parts, fmt.Sprintf(", %s %s", util.CamelCase(a.Name), goType(a.Type)))
	}
	return strings.Join(parts, "")
}

func goType(t ir.Type) string {
	prefix := ""
	if t.IsList {
		prefix = "[]"
	}

	switch t.Kind {
	case ir.KindBase:
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
	case ir.KindEnum:
		return prefix + util.PascalCase(t.Name)
	case ir.KindStruct:
		return prefix + "*" + util.PascalCase(t.Name)
	}

	return prefix + util.PascalCase(t.Name)
}

func goZeroValue(t ir.Type) string {
	if t.IsList {
		return "nil"
	}
	if t.Kind == ir.KindStruct || t.Kind == ir.KindNil {
		return "nil"
	}
	if t.Kind == ir.KindEnum {
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

func goTagString(goTag string, f ir.Field) string {
	if strings.TrimSpace(goTag) == "" {
		return ""
	}

	keys := strings.Split(goTag, ",")
	val := f.Tag
	if strings.TrimSpace(val) == "" {
		val = util.SnakeCase(f.Name)
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
