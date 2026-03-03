# API Documentation

{{if .Note}}{{.Note}}
{{end}}

## API List

| Name | Arguments | Returns | Description |
| :--- | :--- | :--- | :--- |
{{- range .Apis}}
| {{.Name | SnakeCase}} | {{range .Args}}{{.Name}} {{if .Type.IsList}}[{{end}}{{.Type.Name}}{{if .Type.IsList}}]{{end}}<br>{{end}} | {{if ne .Result.Name "nil"}}{{if .Result.IsList}}[{{end}}{{.Result.Name}}{{if .Result.IsList}}]{{end}}{{else}}Void{{end}} | {{.Note}} |
{{- end}}

## RPC Error Codes (HTTP Status)

| Code | Name | Description |
| :--- | :--- | :--- |
| 0 | NoConn | 无法连接 (本地或远程网络故障) |
| 200 | Ok | 请求成功 |
| 400 | ReqErr | 请求错误 (参数序列化失败) |
| 401 | NotAuth | 未授权 (登录失效) |
| 404 | NotExist | 资源不存在 |
| 408 | Timeout | 请求超时 (含重试耗尽) |
| 500 | RespErr | 响应处理错误 (反序列化失败) |

## Usage Demos

### Go Client
```go
import (
    "context"
    "fmt"
    "your_project/sb"
)

func main() {
    client := sb.NewClient("http://localhost:8080")
    client.Retries = 3 // 默认已是 3 次
    
    // Example call
    {{- if .Apis}}
    {{- $api := index .Apis 0}}
    {{- $hasRet := ne $api.Result.Name "nil"}}
    {{if $hasRet}}res, status{{else}}status{{end}} := client.{{$api.Name | PascalCase}}(context.Background() {{range $api.Args}}, {{GoValue .Type.Name}}{{end}})
    
    if status != sb.RpcOk {
        fmt.Printf("Request failed with status: %d\n", status)
        return
    }
    {{if $hasRet}}fmt.Printf("Result: %+v\n", res){{end}}
    {{- else}}
    // No APIs defined
    {{- end}}
}
```

### Go Server
```go
import (
    "net/http"
    "your_project/sb"
)

func main() {
    mux := http.NewServeMux()
    
    // Register API handlers (default middleware is optional)
    {{- range $module, $pkgApis := .Groups}}
    sb.Register{{$module | PascalCase}}(mux)
    {{- end}}
    {{- if not .Groups}}
    sb.RegisterApi(mux) 
    {{- end}}

    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", mux)
}
```

### TypeScript Client
```typescript
import * as sb from "./sb";

async function demo() {
    const client = new sb.RpcClient({
        host: "http://localhost:8080",
        timeout: 5000,
        retries: 3 // 默认已是 3 次
    });

    {{- if .Apis}}
    {{- $api := index .Apis 0}}
    {{- $hasRet := ne $api.Result.Name "nil"}}
    // Example: {{$api.Note}}
    {{if $hasRet}}const [res, status]{{else}}const status{{end}} = await client.{{$api.Name | CamelCase}}({{range $i, $arg := $api.Args}}{{if $i}}, {{end}}{{TsValue .Type.Name}}{{end}});
    
    if (status !== sb.RpcErrCode.Ok) {
        console.error("Request failed with status:", status);
        return;
    }
    {{if $hasRet}}console.log("Data:", res);{{end}}
    {{- else}}
    // No APIs defined
    {{- end}}
}
```

## Types

### Enums

{{- range .Enums}}
#### {{.Name}}
{{if .Note}}> {{.Note}}{{end}}

| ID | Name | Description |
| :--- | :--- | :--- |
{{- range .Children}}
| {{.ID}} | {{.Name}} | {{.Note}} |
{{- end}}

{{- end}}


### Structs

{{- range .Structs}}
#### {{.Name}}
{{if .Note}}> {{.Note}}{{end}}

| Field | Type | Description |
| :--- | :--- | :--- |
{{- range .Fields}}
| {{.Name}} | {{if .Type.IsList}}[{{end}}{{.Type.Name}}{{if .Type.IsList}}]{{end}} | {{.Note}} |
{{- end}}

{{- end}}