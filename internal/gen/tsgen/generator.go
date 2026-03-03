package tsgen

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

// Generate 生成 TypeScript 代码。
func Generate(schema *ir.Schema, cfg gen.Config) error {
	targetDir := filepath.Join(cfg.TsDir, "sb")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("create ts dir %s: %w", targetDir, err)
	}

	if err := generateTypeFile(targetDir); err != nil {
		return err
	}
	if err := generateEnumFile(targetDir, schema.Enums); err != nil {
		return err
	}
	if err := generateStructFile(targetDir, schema.Structs); err != nil {
		return err
	}
	if err := generateRPCFile(targetDir, schema.APIs); err != nil {
		return err
	}
	if err := generateIndexFile(targetDir); err != nil {
		return err
	}
	return nil
}

func generateTypeFile(targetDir string) error {
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
	return os.WriteFile(filepath.Join(targetDir, "type.ts"), []byte(content), 0644)
}

func generateEnumFile(targetDir string, enums []ir.Enum) error {
	var b bytes.Buffer
	if len(enums) == 0 {
		b.WriteString("// No enums defined.\n")
		return os.WriteFile(filepath.Join(targetDir, "enum.ts"), b.Bytes(), 0644)
	}

	for _, e := range enums {
		name := util.PascalCase(e.Name)
		if e.Note != "" {
			fmt.Fprintf(&b, "// %s\n", e.Note)
		}
		fmt.Fprintf(&b, "export enum %s {\n", name)
		for _, m := range e.Members {
			member := util.PascalCase(m.Name)
			if m.Note == "" {
				fmt.Fprintf(&b, "  %s = %d,\n", member, m.ID)
				continue
			}
			fmt.Fprintf(&b, "  %s = %d, // %s\n", member, m.ID, m.Note)
		}
		b.WriteString("}\n\n")
	}

	return os.WriteFile(filepath.Join(targetDir, "enum.ts"), b.Bytes(), 0644)
}

func generateStructFile(targetDir string, structs []ir.Struct) error {
	var b bytes.Buffer
	b.WriteString("import * as Enum from \"./enum\"\n\n")

	for _, st := range structs {
		name := util.PascalCase(st.Name)
		if st.Note != "" {
			fmt.Fprintf(&b, "// %s\n", st.Note)
		}
		fmt.Fprintf(&b, "export interface %s {\n", name)
		for _, f := range st.Fields {
			fmt.Fprintf(&b, "  %s: %s;\n", util.CamelCase(f.Name), tsType(f.Type))
		}
		b.WriteString("}\n\n")

		fmt.Fprintf(&b, "export const new%s = (): %s => ({\n", name, name)
		for _, f := range st.Fields {
			fmt.Fprintf(&b, "  %s: %s,\n", util.CamelCase(f.Name), tsZeroValue(f.Type))
		}
		b.WriteString("})\n\n")
	}

	return os.WriteFile(filepath.Join(targetDir, "struct.ts"), b.Bytes(), 0644)
}

func generateRPCFile(targetDir string, apis []ir.API) error {
	var b bytes.Buffer
	b.WriteString("import { RpcErrCode } from \"./type\"\n")
	b.WriteString("import * as Enum from \"./enum\"\n")
	b.WriteString("import * as Struct from \"./struct\"\n\n")
	b.WriteString("export class RpcClient {\n")
	b.WriteString("  public constructor(public readonly host: string) {}\n\n")

	for _, api := range apis {
		method := util.CamelCase(api.Name)
		args, argNames := tsRPCArgs(api.Args)
		_ = argNames

		if api.Note != "" {
			fmt.Fprintf(&b, "  /** %s */\n", api.Note)
		}
		if api.Result.Kind == ir.KindNil {
			fmt.Fprintf(&b, "  public %s = async (%s): Promise<RpcErrCode> => {\n", method, args)
			b.WriteString("    return RpcErrCode.RespErr\n")
			b.WriteString("  }\n\n")
			continue
		}

		retType := tsRPCType(api.Result)
		fmt.Fprintf(&b, "  public %s = async (%s): Promise<[%s, RpcErrCode]> => {\n", method, args, retType)
		fmt.Fprintf(&b, "    return [%s, RpcErrCode.RespErr]\n", tsRPCZeroValue(api.Result))
		b.WriteString("  }\n\n")
	}
	b.WriteString("}\n\n")

	groups := map[string][]string{}
	for _, api := range apis {
		groups[util.APIGroup(api.Name)] = append(groups[util.APIGroup(api.Name)], api.Name)
	}
	if len(groups) > 0 {
		keys := make([]string, 0, len(groups))
		for k := range groups {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sort.Strings(groups[k])
			fmt.Fprintf(&b, "export const API_%s = [\n", strings.ToUpper(util.SnakeCase(k)))
			for _, n := range groups[k] {
				fmt.Fprintf(&b, "  \"%s\",\n", n)
			}
			b.WriteString("] as const\n\n")
		}
	}

	return os.WriteFile(filepath.Join(targetDir, "rpc.ts"), b.Bytes(), 0644)
}

func generateIndexFile(targetDir string) error {
	content := "export * from \"./type\"\nexport * from \"./enum\"\nexport * from \"./struct\"\nexport * from \"./rpc\"\n"
	return os.WriteFile(filepath.Join(targetDir, "_.ts"), []byte(content), 0644)
}

func tsRPCArgs(args []ir.APIArg) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}
	parts := make([]string, 0, len(args))
	names := make([]string, 0, len(args))
	for _, a := range args {
		name := util.CamelCase(a.Name)
		parts = append(parts, fmt.Sprintf("%s: %s", name, tsRPCType(a.Type)))
		names = append(names, name)
	}
	return strings.Join(parts, ", "), names
}

func tsType(t ir.Type) string {
	base := "unknown"
	switch t.Kind {
	case ir.KindBase:
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
	case ir.KindEnum:
		base = "Enum." + util.PascalCase(t.Name)
	case ir.KindStruct:
		base = util.PascalCase(t.Name)
	case ir.KindNil:
		base = "void"
	}

	if t.IsList {
		return base + "[]"
	}
	return base
}

func tsZeroValue(t ir.Type) string {
	if t.IsList {
		return "[]"
	}
	if t.Kind == ir.KindStruct {
		return "{} as " + util.PascalCase(t.Name)
	}
	if t.Kind == ir.KindEnum {
		return "0 as Enum." + util.PascalCase(t.Name)
	}
	if t.Kind == ir.KindNil {
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

func tsRPCType(t ir.Type) string {
	base := tsType(t)
	if t.Kind == ir.KindStruct {
		base = "Struct." + base
	}
	return base
}

func tsRPCZeroValue(t ir.Type) string {
	if t.IsList {
		return "[]"
	}
	if t.Kind == ir.KindStruct {
		return "Struct.new" + util.PascalCase(t.Name) + "()"
	}
	return tsZeroValue(t)
}
