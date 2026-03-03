

// 账户状态
export enum AccountStatus {
    Offline = 0, 
    Online = 1, 
    Deleted = 2, 
}

export const IsAccountStatus = (v: AccountStatus): boolean => {
    switch (v) {
    case AccountStatus.Offline, AccountStatus.Online, AccountStatus.Deleted:
        return true;
    default:
        return false;
    }
}

export const IsAccountStatusList = (v: AccountStatus[]): boolean => {
    for (const item of v) {
        if (!IsAccountStatus(item)) return false;
    }
    return true;
}


// 类型
export enum Type {
    Sim = 0, 
    Recharge = 1, 
}

export const IsType = (v: Type): boolean => {
    switch (v) {
    case Type.Sim, Type.Recharge:
        return true;
    default:
        return false;
    }
}

export const IsTypeList = (v: Type[]): boolean => {
    for (const item of v) {
        if (!IsType(item)) return false;
    }
    return true;
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

export const IsStatus = (v: Status): boolean => {
    switch (v) {
    case Status.Ok, Status.Err, Status.Two, Status.Three, Status.Four, Status.Five, Status.Six, Status.Seven, Status.One:
        return true;
    default:
        return false;
    }
}

export const IsStatusList = (v: Status[]): boolean => {
    for (const item of v) {
        if (!IsStatus(item)) return false;
    }
    return true;
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

export const IsStatusA = (v: StatusA): boolean => {
    switch (v) {
    case StatusA.Ok, StatusA.One, StatusA.Two, StatusA.Three, StatusA.Four, StatusA.Five, StatusA.Six, StatusA.Seven:
        return true;
    default:
        return false;
    }
}

export const IsStatusAList = (v: StatusA[]): boolean => {
    for (const item of v) {
        if (!IsStatusA(item)) return false;
    }
    return true;
}


// 订单状态
export enum ItemStatus {
    Offline = 0, 
    Online = 1, 
}

export const IsItemStatus = (v: ItemStatus): boolean => {
    switch (v) {
    case ItemStatus.Offline, ItemStatus.Online:
        return true;
    default:
        return false;
    }
}

export const IsItemStatusList = (v: ItemStatus[]): boolean => {
    for (const item of v) {
        if (!IsItemStatus(item)) return false;
    }
    return true;
}


// 可否选号
export enum SimPickPhone {
    No = 0, 
    Yes = 1, 
    Active = 3, 
    Abcc = 4, 
}

export const IsSimPickPhone = (v: SimPickPhone): boolean => {
    switch (v) {
    case SimPickPhone.No, SimPickPhone.Yes, SimPickPhone.Active, SimPickPhone.Abcc:
        return true;
    default:
        return false;
    }
}

export const IsSimPickPhoneList = (v: SimPickPhone[]): boolean => {
    for (const item of v) {
        if (!IsSimPickPhone(item)) return false;
    }
    return true;
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

export const IsSimOperator = (v: SimOperator): boolean => {
    switch (v) {
    case SimOperator.Zz, SimOperator.Lt, SimOperator.Yd, SimOperator.Dx, SimOperator.Gd, SimOperator.Xx, SimOperator.A, SimOperator.B:
        return true;
    default:
        return false;
    }
}

export const IsSimOperatorList = (v: SimOperator[]): boolean => {
    for (const item of v) {
        if (!IsSimOperator(item)) return false;
    }
    return true;
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

export const IsOrderStatus = (v: OrderStatus): boolean => {
    switch (v) {
    case OrderStatus.Pending, OrderStatus.Closed, OrderStatus.Canceled, OrderStatus.Shipped, OrderStatus.Delivered, OrderStatus.Actived, OrderStatus.Settled:
        return true;
    default:
        return false;
    }
}

export const IsOrderStatusList = (v: OrderStatus[]): boolean => {
    for (const item of v) {
        if (!IsOrderStatus(item)) return false;
    }
    return true;
}

