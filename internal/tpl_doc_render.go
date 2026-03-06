package internal

import (
	"fmt"
	"strings"
)

func generateDoc(schema *TplSchema, cfg Config) error {
	content := renderDoc(schema)
	if err := writeDocFile(cfg.GoDir, []byte(content)); err != nil {
		return err
	}
	if err := writeDocFile(cfg.TsDir, []byte(content)); err != nil {
		return err
	}
	return nil
}

func renderDoc(schema *TplSchema) string {
	var w sourceWriter
	w.Line("# API Documentation")
	w.Blank()
	if schema.Note != "" {
		w.Line(schema.Note)
		w.Blank()
	}
	w.Line("## API List")
	w.Blank()
	w.Line("| Name | Arguments | Returns | Description |")
	w.Line("| :--- | :--- | :--- | :--- |")
	for _, api := range schema.Apis {
		w.Linef("| %s | %s | %s | %s |", SnakeCase(api.Name), renderDocArgs(api.Args), renderDocReturn(api.Result), tplRenderMarkdownInline(api.Note))
	}
	w.Blank()
	w.Line("## RPC Error Codes (HTTP Status)")
	w.Blank()
	w.Line("| Code | Name | Description |")
	w.Line("| :--- | :--- | :--- |")
	w.Line("| 0 | NoConn | 无法连接 (本地或远程网络故障) |")
	w.Line("| 200 | Ok | 请求成功 |")
	w.Line("| 400 | ReqErr | 请求错误 (参数序列化失败) |")
	w.Line("| 401 | NotAuth | 未授权 (登录失效) |")
	w.Line("| 404 | NotExist | 资源不存在 |")
	w.Line("| 408 | Timeout | 请求超时 (含重试耗尽) |")
	w.Line("| 500 | RespErr | 响应处理错误 (反序列化失败) |")
	w.Blank()
	w.Line("## Usage Demos")
	w.Blank()
	w.Line("### Go Client")
	w.Line("```go")
	w.Line("import (")
	w.Line("    \"context\"")
	w.Line("    \"fmt\"")
	w.Line("    \"your_project/sb\"")
	w.Line(")")
	w.Blank()
	w.Line("func main() {")
	w.Line("    client := sb.NewClient(\"http://localhost:8080\")")
	w.Line("    client.Retries = 3 // 默认已是 3 次")
	w.Line("    ")
	w.Line("    // Example call")
	if len(schema.Apis) > 0 {
		api := schema.Apis[0]
		args := renderGoExampleArgs(api.Args)
		if api.Result.Name == "nil" {
			w.Linef("    status := sb.Call%s(client, context.Background()%s)", PascalCase(api.Name), args)
		} else {
			w.Linef("    res, status := sb.Call%s(client, context.Background()%s)", PascalCase(api.Name), args)
		}
		w.Line("    ")
		w.Line("    if status != sb.RpcOk {")
		w.Line("        fmt.Printf(\"Request failed with status: %d\\n\", status)")
		w.Line("        return")
		w.Line("    }")
		if api.Result.Name != "nil" {
			w.Line("    fmt.Printf(\"Result: %+v\\n\", res)")
		}
	} else {
		w.Line("    // No APIs defined")
	}
	w.Line("}")
	w.Line("```")
	w.Blank()
	w.Line("### Go Server")
	w.Line("```go")
	w.Line("import (")
	w.Line("    \"fmt\"")
	w.Line("    \"net/http\"")
	w.Line("    \"time\"")
	w.Line("    \"your_project/sb\"")
	w.Line(")")
	w.Blank()
	w.Line("func main() {")
	w.Line("    mux := http.NewServeMux()")
	w.Line("    ")
	w.Line("    // Request-level timeout propagates to business logic and downstream calls.")
	w.Line("    rpcTimeout := 3 * time.Second")
	w.Line("    ")
	w.Line("    // Connection-level timeouts protect the HTTP transport.")
	w.Line("    server := &http.Server{")
	w.Line("        Addr:              \":8080\",")
	w.Line("        Handler:           mux,")
	w.Line("        ReadHeaderTimeout: 2 * time.Second,")
	w.Line("        ReadTimeout:       5 * time.Second,")
	w.Line("        WriteTimeout:      8 * time.Second,")
	w.Line("        IdleTimeout:       30 * time.Second,")
	w.Line("    }")
	w.Line("    ")
	w.Line("    // Register API handlers (middleware is optional)")
	groups := groupAPIs(schema.Apis)
	keys := orderedGroupKeys(groups)
	if len(keys) == 0 {
		w.Line("    sb.RegisterApi(mux, sb.TimeoutMiddleware(rpcTimeout))")
	} else {
		for _, key := range keys {
			funcName := "Register" + PascalCase(key) + "Api"
			if PascalCase(key) == "Api" {
				funcName = "RegisterApi"
			}
			w.Linef("    sb.%s(mux, sb.TimeoutMiddleware(rpcTimeout))", funcName)
		}
	}
	w.Blank()
	w.Line("    fmt.Println(\"Server starting on\", server.Addr)")
	w.Line("    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {")
	w.Line("        panic(err)")
	w.Line("    }")
	w.Line("}")
	w.Line("```")
	w.Blank()
	w.Line("### TypeScript Client")
	w.Line("```typescript")
	w.Line("import * as sb from \"./sb\";")
	w.Blank()
	w.Line("async function demo() {")
	w.Line("    const client = new sb.RpcClient({")
	w.Line("        host: \"http://localhost:8080\",")
	w.Line("        timeout: 5000,")
	w.Line("        retries: 3 // 默认已是 3 次")
	w.Line("    });")
	if len(schema.Apis) > 0 {
		api := schema.Apis[0]
		if api.Note != "" {
			w.WriteLineCommentWithHead("    // ", "Example: ", api.Note)
			w.Blank()
		}
		args := renderTsExampleArgs(api.Args)
		if api.Result.Name == "nil" {
			w.Linef("    const status = await client.%s(%s);", CamelCase(api.Name), args)
		} else {
			w.Linef("    const [res, status] = await client.%s(%s);", CamelCase(api.Name), args)
		}
		w.Line("    ")
		w.Line("    if (status !== sb.RpcErrCode.Ok) {")
		w.Line("        console.error(\"Request failed with status:\", status);")
		w.Line("        return;")
		w.Line("    }")
		if api.Result.Name != "nil" {
			w.Line("    console.log(\"Data:\", res);")
		}
	} else {
		w.Line("    // No APIs defined")
	}
	w.Line("}")
	w.Line("```")
	w.Blank()
	w.Line("## Types")
	w.Blank()
	w.Line("### Enums")
	for _, enum := range schema.Enums {
		w.Linef("#### %s", enum.Name)
		if enum.Note != "" {
			w.Line(tplRenderMarkdownQuote(enum.Note))
			w.Blank()
		}
		w.Line("| ID | Name | Description |")
		w.Line("| :--- | :--- | :--- |")
		for _, child := range enum.Children {
			w.Linef("| %d | %s | %s |", child.ID, PascalCase(child.Name), tplRenderMarkdownInline(child.Note))
		}
		w.Blank()
	}
	w.Blank()
	w.Line("### Structs")
	for _, st := range schema.Structs {
		w.Linef("#### %s", st.Name)
		if st.Note != "" {
			w.Line(tplRenderMarkdownQuote(st.Note))
			w.Blank()
		}
		w.Line("| Field | Type | Description |")
		w.Line("| :--- | :--- | :--- |")
		for _, field := range st.Fields {
			fieldType := field.Type.Name
			if field.Type.IsList {
				fieldType = "[" + fieldType + "]"
			}
			w.Linef("| %s | %s | %s |", field.Name, fieldType, tplRenderMarkdownInline(field.Note))
		}
		w.Blank()
	}
	return w.String()
}

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
