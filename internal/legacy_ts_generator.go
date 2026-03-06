package internal

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// Generate 生成 TypeScript 代码。
func GenerateLegacyTs(schema *IRSchema, cfg Config) error {
	dir, err := newGeneratedDir(filepath.Join(cfg.TsDir, "sb"))
	if err != nil {
		return err
	}

	if err := legacyTsGenerateTypeFile(dir); err != nil {
		return err
	}
	if err := legacyTsGenerateEnumFile(dir, schema.Enums); err != nil {
		return err
	}
	if err := legacyTsGenerateStructFile(dir, schema.Structs); err != nil {
		return err
	}
	if err := legacyTsGenerateRPCFile(dir, schema.APIs); err != nil {
		return err
	}
	if err := legacyTsGenerateIndexFile(dir); err != nil {
		return err
	}
	return nil
}

func legacyTsGenerateTypeFile(dir *generatedDir) error {
	content := `export enum RpcErrCode {
	Ok = 200,
	NoConn = 0,
	Timeout = 408,
	ReqErr = 400,
	RespErr = 500,
	NotAuth = 401,
	NotExist = 404,
}
`
	return dir.Write("type.ts", []byte(content), 0644)
}

func legacyTsGenerateEnumFile(dir *generatedDir, enums []IREnum) error {
	var b bytes.Buffer
	if len(enums) == 0 {
		b.WriteString("// No enums defined.\n")
		return dir.Write("enum.ts", b.Bytes(), 0644)
	}

	for _, e := range enums {
		name := PascalCase(e.Name)
		if e.Note != "" {
			b.WriteString(RenderLineComment("// ", e.Note))
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "export enum %s {\n", name)
		for _, m := range e.Members {
			member := PascalCase(m.Name)
			if m.Note != "" {
				b.WriteString(RenderLineComment("  // ", m.Note))
				b.WriteByte('\n')
			}
			fmt.Fprintf(&b, "  %s = %d,\n", member, m.ID)
		}
		b.WriteString("}\n\n")
	}

	return dir.Write("enum.ts", b.Bytes(), 0644)
}

func legacyTsGenerateStructFile(dir *generatedDir, structs []IRStruct) error {
	var b bytes.Buffer
	b.WriteString("import * as Enum from \"./enum\"\n\n")

	for _, st := range structs {
		name := PascalCase(st.Name)
		if st.Note != "" {
			b.WriteString(RenderLineComment("// ", st.Note))
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "export interface %s {\n", name)
		for _, f := range st.Fields {
			fmt.Fprintf(&b, "  %s: %s;\n", CamelCase(f.Name), legacyTsType(f.Type))
		}
		b.WriteString("}\n\n")

		fmt.Fprintf(&b, "export const new%s = (): %s => ({\n", name, name)
		for _, f := range st.Fields {
			fmt.Fprintf(&b, "  %s: %s,\n", CamelCase(f.Name), legacyTsZeroValue(f.Type))
		}
		b.WriteString("})\n\n")
	}

	return dir.Write("struct.ts", b.Bytes(), 0644)
}

func legacyTsGenerateRPCFile(dir *generatedDir, apis []IRAPI) error {
	var b bytes.Buffer
	b.WriteString("import { RpcErrCode } from \"./type\"\n")
	b.WriteString("import * as Enum from \"./enum\"\n")
	b.WriteString("import * as Struct from \"./struct\"\n\n")
	b.WriteString("export class RpcClient {\n")
	b.WriteString("  public constructor(public readonly host: string) {}\n\n")

	for _, api := range apis {
		method := CamelCase(api.Name)
		args, argNames := legacyTsRPCArgs(api.Args)
		_ = argNames

		if api.Note != "" {
			b.WriteString(RenderLineComment("  // ", api.Note))
			b.WriteByte('\n')
		}
		if api.Result.Kind == IRKindNil {
			fmt.Fprintf(&b, "  public %s = async (%s): Promise<RpcErrCode> => {\n", method, args)
			b.WriteString("    return RpcErrCode.RespErr\n")
			b.WriteString("  }\n\n")
			continue
		}

		retType := legacyTsRPCType(api.Result)
		fmt.Fprintf(&b, "  public %s = async (%s): Promise<[%s, RpcErrCode]> => {\n", method, args, retType)
		fmt.Fprintf(&b, "    return [%s, RpcErrCode.RespErr]\n", legacyTsRPCZeroValue(api.Result))
		b.WriteString("  }\n\n")
	}
	b.WriteString("}\n\n")

	groups := map[string][]string{}
	for _, api := range apis {
		groups[APIGroup(api.Name)] = append(groups[APIGroup(api.Name)], api.Name)
	}
	if len(groups) > 0 {
		keys := make([]string, 0, len(groups))
		for k := range groups {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sort.Strings(groups[k])
			fmt.Fprintf(&b, "export const API_%s = [\n", strings.ToUpper(SnakeCase(k)))
			for _, n := range groups[k] {
				fmt.Fprintf(&b, "  \"%s\",\n", n)
			}
			b.WriteString("] as const\n\n")
		}
	}

	return dir.Write("rpc.ts", b.Bytes(), 0644)
}

func legacyTsGenerateIndexFile(dir *generatedDir) error {
	content := "export * from \"./type\"\nexport * from \"./enum\"\nexport * from \"./struct\"\nexport * from \"./rpc\"\n"
	return dir.Write("_.ts", []byte(content), 0644)
}

func legacyTsRPCArgs(args []IRAPIArg) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}
	parts := make([]string, 0, len(args))
	names := make([]string, 0, len(args))
	for _, a := range args {
		name := CamelCase(a.Name)
		parts = append(parts, fmt.Sprintf("%s: %s", name, legacyTsRPCType(a.Type)))
		names = append(names, name)
	}
	return strings.Join(parts, ", "), names
}

func legacyTsType(t IRType) string {
	base := "unknown"
	switch t.Kind {
	case IRKindBase:
		switch t.Name {
		case "i8", "u8", "i16", "u16", "i32", "u32", "f32", "f64":
			base = "number"
		case "i64", "u64":
			base = "bigint"
		case "bool":
			base = "boolean"
		case "text":
			base = "string"
		case "bin":
			base = "Uint8Array"
		}
	case IRKindEnum:
		base = "Enum." + PascalCase(t.Name)
	case IRKindStruct:
		base = PascalCase(t.Name)
	case IRKindNil:
		base = "void"
	}

	if t.IsList {
		return base + "[]"
	}
	return base
}

func legacyTsZeroValue(t IRType) string {
	if t.IsList {
		return "[]"
	}
	if t.Kind == IRKindStruct {
		return "{} as " + PascalCase(t.Name)
	}
	if t.Kind == IRKindEnum {
		return "0 as Enum." + PascalCase(t.Name)
	}
	if t.Kind == IRKindNil {
		return "undefined"
	}

	switch t.Name {
	case "bool":
		return "false"
	case "text":
		return "\"\""
	case "bin":
		return "new Uint8Array(0)"
	case "i64", "u64":
		return "0n"
	default:
		return "0"
	}
}

func legacyTsRPCType(t IRType) string {
	base := legacyTsType(t)
	if t.Kind == IRKindStruct {
		base = "Struct." + base
	}
	return base
}

func legacyTsRPCZeroValue(t IRType) string {
	if t.IsList {
		return "[]"
	}
	if t.Kind == IRKindStruct {
		return "Struct.new" + PascalCase(t.Name) + "()"
	}
	return legacyTsZeroValue(t)
}
