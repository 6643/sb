
// 账户状态
export enum AccountStatus {
    Offline = 0, 
    Online = 1, 
    Deleted = 2, 
}

// 类型
export enum Type {
    Sim = 0, 
    Recharge = 1, 
}

// 错误码
export enum Status {
    Ok = 0, 
    Err = 1, 
    Two = 2, 
    Three = 3, 
    Four = 4, 
    Five = 5, 
    Six = 6, 
    Seven = 7, 
    One = 11, 
}

// 状态A
export enum StatusA {
    Ok = 0, 
    One = 1, 
    Two = 2, 
    Three = 3, 
    Four = 4, 
    Five = 5, 
    Six = 6, 
    Seven = 7, 
}

// 订单状态
export enum ItemStatus {
    Offline = 0, 
    Online = 1, 
}

// 可否选号
export enum SimPickPhone {
    No = 0, 
    Yes = 1, 
    Active = 3, 
    Abcc = 4, 
}

// 运营商
export enum SimOperator {
    Zz = 2, 
    Lt = 3, 
    Yd = 4, 
    Dx = 5, 
    Gd = 6, 
    Xx = 7, 
    A = 11, 
    B = 12, 
}

// 订单状态
export enum OrderStatus {
    Pending = 0, // 待处理
    Closed = 1, // 已关闭
    Canceled = 2, // 已取消
    Shipped = 3, // 已发货
    Delivered = 4, // 已送达
    Actived = 5, // 已激活
    Settled = 6, // 已结算
}
