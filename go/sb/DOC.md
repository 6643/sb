# API Documentation

## API List

| Name | Arguments | Returns | Description |
| :--- | :--- | :--- | :--- |
| user_get_abc |  | OrderStatus | 获取用户的id |
| user_get_abcd | page u8<br>size u8<br> | OrderStatus | 获取abcd |
| user_set_sim_info | info SimInfo<br> | Void | 设置sim信息<br>无返回值 |
| get_count | page u8<br> | u8 | 获取数量 |
| get_bin | page u8<br> | bin | 获取bin |

## RPC Status Codes

| Code | Name | Description |
| :--- | :--- | :--- |
| 0 | NoConn | 客户端本地错误或网络不可达; 服务端不会写出这个 HTTP 状态码 |
| 200 | Ok | 请求成功 |
| 400 | ReqErr | 请求错误 (参数序列化失败) |
| 401 | NotAuth | 未授权 (登录失效) |
| 404 | NotExist | 资源不存在 |
| 408 | Timeout | 单次请求在发送或读取响应体阶段超时, 或服务端超时 |
| 500 | RespErr | 响应处理错误 (反序列化失败、响应体损坏或内部错误) |

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
    client.Timeout = 5 * time.Second // 单次请求覆盖建连到响应体读取
    client.Retries = 3 // 默认已是 3 次
    
    // Example call
    res, status := sb.CallUserGetAbc(client, context.Background())
    
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
    "fmt"
    "net/http"
    "time"
    "your_project/sb"
)

func main() {
    mux := http.NewServeMux()
    
    // Request-level timeout propagates to business logic and downstream calls.
    rpcTimeout := 3 * time.Second
    
    // Connection-level timeouts protect the HTTP transport.
    server := &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 2 * time.Second,
        ReadTimeout:       5 * time.Second,
        WriteTimeout:      8 * time.Second,
        IdleTimeout:       30 * time.Second,
    }
    
    // Register API handlers (middleware is optional)
    sb.RegisterApi(mux, sb.TimeoutMiddleware(rpcTimeout))
    sb.RegisterUserApi(mux, sb.TimeoutMiddleware(rpcTimeout))

    fmt.Println("Server starting on", server.Addr)
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic(err)
    }
}
```

### TypeScript Client
```typescript
import * as sb from "./sb";

async function demo() {
    const client = new sb.RpcClient({
        host: "http://localhost:8080",
        timeout: 5000,
        retries: 3, // 默认已是 3 次
        maxRespBytes: 4 * 1024 * 1024, // 传输中超限会立即失败
    });
    // Example: 获取用户的id

    const [res, status] = await client.userGetAbc();
    
    if (status !== sb.RpcErrCode.Ok) {
        console.error("Request failed with status:", status); // 未知 HTTP 状态可能直接返回数字
        return;
    }
    console.log("Data:", res);
}
```

## Types

### Enums
#### AccountStatus
> 账户状态

| ID | Name | Description |
| :--- | :--- | :--- |
| 0 | Offline |  |
| 1 | Online |  |
| 2 | Deleted |  |

#### Type
> 类型

| ID | Name | Description |
| :--- | :--- | :--- |
| 0 | Sim |  |
| 1 | Recharge |  |

#### Status
> 错误码

| ID | Name | Description |
| :--- | :--- | :--- |
| 0 | Ok |  |
| 1 | Err |  |
| 2 | Two |  |
| 3 | Three |  |
| 4 | Four |  |
| 5 | Five |  |
| 6 | Six |  |
| 7 | Seven |  |
| 11 | One |  |

#### StatusA
> 状态A

| ID | Name | Description |
| :--- | :--- | :--- |
| 0 | Ok |  |
| 1 | One |  |
| 2 | Two |  |
| 3 | Three |  |
| 4 | Four |  |
| 5 | Five |  |
| 6 | Six |  |
| 7 | Seven |  |

#### ItemStatus
> 订单状态

| ID | Name | Description |
| :--- | :--- | :--- |
| 0 | Offline |  |
| 1 | Online |  |

#### SimPickPhone
> 可否选号

| ID | Name | Description |
| :--- | :--- | :--- |
| 0 | No |  |
| 1 | Yes |  |
| 3 | Active |  |
| 4 | Abcc |  |

#### SimOperator
> 运营商

| ID | Name | Description |
| :--- | :--- | :--- |
| 2 | Zz |  |
| 3 | Lt |  |
| 4 | Yd |  |
| 5 | Dx |  |
| 6 | Gd |  |
| 7 | Xx |  |
| 11 | A |  |
| 12 | B |  |

#### OrderStatus
> 订单状态
> Pending      待处理
> Closed       已关闭
> Canceled     已取消
> Shipped      已发货
> Delivered    已送达
> Actived      已激活
> Settled      已结算

| ID | Name | Description |
| :--- | :--- | :--- |
| 0 | Pending |  |
| 1 | Closed |  |
| 2 | Canceled |  |
| 3 | Shipped |  |
| 4 | Delivered |  |
| 5 | Actived |  |
| 6 | Settled |  |


### Structs
#### Recharge
| Field | Type | Description |
| :--- | :--- | :--- |
| id | u32 | abcd |
| type | [OrderStatus] |  |
| phone | [text] |  |
| si | SimInfo |  |

#### RechargeA
| Field | Type | Description |
| :--- | :--- | :--- |
| id | u32 | abcd |
| type | [OrderStatus] |  |
| phone | [text] |  |
| si | SimInfo |  |
| aid | u32 |  |

#### RechargeB
| Field | Type | Description |
| :--- | :--- | :--- |
| id | u32 | abcd |
| type | [OrderStatus] |  |
| phone | [text] |  |
| si | SimInfo |  |
| bid | u32 |  |

#### Sim
| Field | Type | Description |
| :--- | :--- | :--- |
| id | u32 | SIM卡ID |
| type | Type |  |
| status | ItemStatus |  |
| commission | u16 | 佣金 |
| supplier | u32 | 供应商ID |
| aff | u32 | 推广员ID |
| contract_duration | u8 | 合约期(月), 0:长期 |
| name | text |  |
| operator | SimOperator | 运营商 |
| monthly | u16 | 月租 |
| flow_universal | u16 | 通用流量 |
| flow_directional | u16 | 定向流量 |
| can_move_flow | bool | 流量是否结转 |
| call_month | u16 | 每月通话(分钟) |
| call_price | u16 |  |
| sms_month | u16 | 每月短信(条) |
| sms_price | u16 |  |
| min_age | u8 |  |
| max_age | u8 |  |
| attribution | u32 | 归属地, 0:随机, 1:收货地 |
| pick_phone | [SimPickPhone] | 选号 |
| first_charge_link | text | 首充渠道 |
| first_charge_money | text | 首充金额 |
| first_charge_return | text | 首充返额 |
| ban_city | [u32] | 禁发区域 |
| info | [SimInfo] |  |
| snapshot | [text] | 套餐截图 |

#### SimInfo
| Field | Type | Description |
| :--- | :--- | :--- |
| id | u32 |  |
| title | text |  |
| content | text |  |
| a | bool |  |
| b | bool |  |
| c | bool |  |
| d | bool |  |
| zip | bin |  |

#### SimOrder2
| Field | Type | Description |
| :--- | :--- | :--- |
| id | u32 | SIM卡ID |
| name | text | 办理人姓名 |
| phone | text | 联系电话 |
| id_no | text | 身份证号 |
| city_code | u32 | 所在城市 |
| address | text | 详细地址 |
| new_phone | text | 新手机号码 |

#### SimOrder
| Field | Type | Description |
| :--- | :--- | :--- |
| id | u32 |  |
| account_id | u32 |  |
| item_id | u32 |  |
| name | text | 办理人姓名 |
| phone | text | 联系电话 |
| id_no | text | 身份证号 |
| city_code | u32 | 所在城市 |
| address | text | 详细地址 |
| new_phone | text | 新手机号码 |
| commission | u16 | 佣金 |
| status | OrderStatus |  |

