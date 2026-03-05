# API Documentation



## API List

| Name | Arguments | Returns | Description |
| :--- | :--- | :--- | :--- |
| echo | env Envelope<br> | Envelope |  |
| echo_tail | env Envelope<br> | Envelope |  |
| get_color |  | Color |  |
| get_bad_color |  | Color |  |
| ping |  | Void |  |
| ping_junk |  | Void |  |
| pick | color Color<br> | Color |  |
| pick_bad | color Color<br> | Color |  |

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
    res, status := sb.CallEcho(client, context.Background() , 0)
    
    if status != sb.RpcOk {
        fmt.Printf("Request failed with status: %d\n", status)
        return
    }
    fmt.Printf("Result: %+v\n", res)
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
    sb.RegisterApi(mux)

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
    // Example: 
    const [res, status] = await client.echo(0);
    
    if (status !== sb.RpcErrCode.Ok) {
        console.error("Request failed with status:", status);
        return;
    }
    console.log("Data:", res);
}
```

## Types

### Enums
#### Color


| ID | Name | Description |
| :--- | :--- | :--- |
| 1 | Red |  |
| 2 | Green |  |
| 3 | Blue |  |


### Structs
#### Item


| Field | Type | Description |
| :--- | :--- | :--- |
| id | u32 |  |
| color | Color |  |
| tags | [text] |  |
| active | bool |  |
#### Envelope


| Field | Type | Description |
| :--- | :--- | :--- |
| item | Item |  |
| items | [Item] |  |
| note | text |  |
