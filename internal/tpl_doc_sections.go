package internal

func renderDocAPIList(w *sourceWriter, apis []TplApi) {
	w.Line("## API List")
	w.Blank()
	w.Line("| Name | Arguments | Returns | Description |")
	w.Line("| :--- | :--- | :--- | :--- |")
	for _, api := range apis {
		w.Linef("| %s | %s | %s | %s |", SnakeCase(api.Name), renderDocArgs(api.Args), renderDocReturn(api.Result), tplRenderMarkdownInline(api.Note))
	}
	w.Blank()
}

func renderDocRPCStatusCodes(w *sourceWriter) {
	w.Line("## RPC Status Codes")
	w.Blank()
	w.Line("| Code | Name | Description |")
	w.Line("| :--- | :--- | :--- |")
	w.Line("| 0 | NoConn | 客户端本地错误或网络不可达; 服务端不会写出这个 HTTP 状态码 |")
	w.Line("| 200 | Ok | 请求成功 |")
	w.Line("| 400 | ReqErr | 请求错误 (参数序列化失败) |")
	w.Line("| 401 | NotAuth | 未授权 (登录失效) |")
	w.Line("| 404 | NotExist | 资源不存在 |")
	w.Line("| 408 | Timeout | 单次请求在发送或读取响应体阶段超时, 或服务端超时 |")
	w.Line("| 500 | RespErr | 响应处理错误 (反序列化失败、响应体损坏或内部错误) |")
	w.Blank()
}

func renderDocUsageDemos(w *sourceWriter, schema *TplSchema) {
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
	w.Line("    client.Timeout = 5 * time.Second // 单次请求覆盖建连到响应体读取")
	w.Line("    client.EnableRetries = false // 默认关闭, 只有请求具备幂等保障时再开启")
	w.Line("    client.Retries = 3 // 仅在 EnableRetries=true 时生效")
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
	w.Line("    // Request-level timeout returns HTTP 408 and propagates cancellation to business logic.")
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
	w.Line("        enableRetries: false, // 默认关闭, 只有请求具备幂等保障时再开启")
	w.Line("        retries: 3, // 仅在 enableRetries=true 时生效")
	w.Line("        maxRespBytes: 4 * 1024 * 1024, // 传输中超限会立即失败")
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
		w.Line("        console.error(\"Request failed with status:\", status); // 未知 HTTP 状态可能直接返回数字")
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
}

func renderDocTypes(w *sourceWriter, schema *TplSchema) {
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
}
