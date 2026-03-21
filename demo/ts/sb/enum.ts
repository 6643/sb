import * as rt from "./type"

// 账户状态
export enum AccountStatus {
    Offline = 0,
    Online = 1,
    Deleted = 2,
}

export const DefaultAccountStatus = (): AccountStatus => AccountStatus.Offline;
export const IsAccountStatus = (v: AccountStatus): boolean => {
    switch (v) {
    case AccountStatus.Offline:
    case AccountStatus.Online:
    case AccountStatus.Deleted:
        return true;
    default:
        return false;
    }
};

export const NormalizeAccountStatus = (v: AccountStatus): AccountStatus => {
    if (IsAccountStatus(v)) return v;
    if ((v as any) === 0) return DefaultAccountStatus();
    return v;
};
export const IsDefaultAccountStatus = (v: AccountStatus): boolean => NormalizeAccountStatus(v) === DefaultAccountStatus();
export const IsAssignableAccountStatus = (v: AccountStatus): boolean => IsAccountStatus(v) || (((v as any) === 0) && !IsAccountStatus(0 as any));
export const eqAccountStatusValue = (a: AccountStatus, b: AccountStatus): boolean => NormalizeAccountStatus(a) === NormalizeAccountStatus(b);
export const eqAccountStatusList = (a: AccountStatus[], b: AccountStatus[]): boolean => rt.eqList(a, b, eqAccountStatusValue);

export const getAccountStatusListBody = (buf: rt.Buffer, state: number): [AccountStatus[], rt.Err] => {
    const [list, err] = rt.getDefaultList<AccountStatus>(
        buf,
        state,
        () => DefaultAccountStatus(),
        (buf) => {
            const [value, err] = rt.getU8(buf);
            if (err !== null) return [DefaultAccountStatus(), rt.errU(err)];
            const item = value as AccountStatus;
            if (!IsAccountStatus(item)) return [DefaultAccountStatus(), new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
    return [list, err];
};

export const setAccountStatusListBody = (buf: rt.Buffer, state: number, v: AccountStatus[]): rt.Err => {
    return rt.setDefaultList<AccountStatus>(
        buf,
        state,
        v,
        (item) => IsDefaultAccountStatus(item),
        (buf, item) => {
            if (!IsAccountStatus(item)) return new Error(`非法枚举值: ${item}`);
            return rt.setU8(buf, item as any);
        },
    );
};

// 类型
export enum Type {
    Sim = 0,
    Recharge = 1,
}

export const DefaultType = (): Type => Type.Sim;
export const IsType = (v: Type): boolean => {
    switch (v) {
    case Type.Sim:
    case Type.Recharge:
        return true;
    default:
        return false;
    }
};

export const NormalizeType = (v: Type): Type => {
    if (IsType(v)) return v;
    if ((v as any) === 0) return DefaultType();
    return v;
};
export const IsDefaultType = (v: Type): boolean => NormalizeType(v) === DefaultType();
export const IsAssignableType = (v: Type): boolean => IsType(v) || (((v as any) === 0) && !IsType(0 as any));
export const eqTypeValue = (a: Type, b: Type): boolean => NormalizeType(a) === NormalizeType(b);
export const eqTypeList = (a: Type[], b: Type[]): boolean => rt.eqList(a, b, eqTypeValue);

export const getTypeListBody = (buf: rt.Buffer, state: number): [Type[], rt.Err] => {
    const [list, err] = rt.getDefaultList<Type>(
        buf,
        state,
        () => DefaultType(),
        (buf) => {
            const [value, err] = rt.getU8(buf);
            if (err !== null) return [DefaultType(), rt.errU(err)];
            const item = value as Type;
            if (!IsType(item)) return [DefaultType(), new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
    return [list, err];
};

export const setTypeListBody = (buf: rt.Buffer, state: number, v: Type[]): rt.Err => {
    return rt.setDefaultList<Type>(
        buf,
        state,
        v,
        (item) => IsDefaultType(item),
        (buf, item) => {
            if (!IsType(item)) return new Error(`非法枚举值: ${item}`);
            return rt.setU8(buf, item as any);
        },
    );
};

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

export const DefaultStatus = (): Status => Status.Ok;
export const IsStatus = (v: Status): boolean => {
    switch (v) {
    case Status.Ok:
    case Status.Err:
    case Status.Two:
    case Status.Three:
    case Status.Four:
    case Status.Five:
    case Status.Six:
    case Status.Seven:
    case Status.One:
        return true;
    default:
        return false;
    }
};

export const NormalizeStatus = (v: Status): Status => {
    if (IsStatus(v)) return v;
    if ((v as any) === 0) return DefaultStatus();
    return v;
};
export const IsDefaultStatus = (v: Status): boolean => NormalizeStatus(v) === DefaultStatus();
export const IsAssignableStatus = (v: Status): boolean => IsStatus(v) || (((v as any) === 0) && !IsStatus(0 as any));
export const eqStatusValue = (a: Status, b: Status): boolean => NormalizeStatus(a) === NormalizeStatus(b);
export const eqStatusList = (a: Status[], b: Status[]): boolean => rt.eqList(a, b, eqStatusValue);

export const getStatusListBody = (buf: rt.Buffer, state: number): [Status[], rt.Err] => {
    const [list, err] = rt.getDefaultList<Status>(
        buf,
        state,
        () => DefaultStatus(),
        (buf) => {
            const [value, err] = rt.getU8(buf);
            if (err !== null) return [DefaultStatus(), rt.errU(err)];
            const item = value as Status;
            if (!IsStatus(item)) return [DefaultStatus(), new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
    return [list, err];
};

export const setStatusListBody = (buf: rt.Buffer, state: number, v: Status[]): rt.Err => {
    return rt.setDefaultList<Status>(
        buf,
        state,
        v,
        (item) => IsDefaultStatus(item),
        (buf, item) => {
            if (!IsStatus(item)) return new Error(`非法枚举值: ${item}`);
            return rt.setU8(buf, item as any);
        },
    );
};

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

export const DefaultStatusA = (): StatusA => StatusA.Ok;
export const IsStatusA = (v: StatusA): boolean => {
    switch (v) {
    case StatusA.Ok:
    case StatusA.One:
    case StatusA.Two:
    case StatusA.Three:
    case StatusA.Four:
    case StatusA.Five:
    case StatusA.Six:
    case StatusA.Seven:
        return true;
    default:
        return false;
    }
};

export const NormalizeStatusA = (v: StatusA): StatusA => {
    if (IsStatusA(v)) return v;
    if ((v as any) === 0) return DefaultStatusA();
    return v;
};
export const IsDefaultStatusA = (v: StatusA): boolean => NormalizeStatusA(v) === DefaultStatusA();
export const IsAssignableStatusA = (v: StatusA): boolean => IsStatusA(v) || (((v as any) === 0) && !IsStatusA(0 as any));
export const eqStatusAValue = (a: StatusA, b: StatusA): boolean => NormalizeStatusA(a) === NormalizeStatusA(b);
export const eqStatusAList = (a: StatusA[], b: StatusA[]): boolean => rt.eqList(a, b, eqStatusAValue);

export const getStatusAListBody = (buf: rt.Buffer, state: number): [StatusA[], rt.Err] => {
    const [list, err] = rt.getDefaultList<StatusA>(
        buf,
        state,
        () => DefaultStatusA(),
        (buf) => {
            const [value, err] = rt.getU8(buf);
            if (err !== null) return [DefaultStatusA(), rt.errU(err)];
            const item = value as StatusA;
            if (!IsStatusA(item)) return [DefaultStatusA(), new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
    return [list, err];
};

export const setStatusAListBody = (buf: rt.Buffer, state: number, v: StatusA[]): rt.Err => {
    return rt.setDefaultList<StatusA>(
        buf,
        state,
        v,
        (item) => IsDefaultStatusA(item),
        (buf, item) => {
            if (!IsStatusA(item)) return new Error(`非法枚举值: ${item}`);
            return rt.setU8(buf, item as any);
        },
    );
};

// 订单状态
export enum ItemStatus {
    Offline = 0,
    Online = 1,
}

export const DefaultItemStatus = (): ItemStatus => ItemStatus.Offline;
export const IsItemStatus = (v: ItemStatus): boolean => {
    switch (v) {
    case ItemStatus.Offline:
    case ItemStatus.Online:
        return true;
    default:
        return false;
    }
};

export const NormalizeItemStatus = (v: ItemStatus): ItemStatus => {
    if (IsItemStatus(v)) return v;
    if ((v as any) === 0) return DefaultItemStatus();
    return v;
};
export const IsDefaultItemStatus = (v: ItemStatus): boolean => NormalizeItemStatus(v) === DefaultItemStatus();
export const IsAssignableItemStatus = (v: ItemStatus): boolean => IsItemStatus(v) || (((v as any) === 0) && !IsItemStatus(0 as any));
export const eqItemStatusValue = (a: ItemStatus, b: ItemStatus): boolean => NormalizeItemStatus(a) === NormalizeItemStatus(b);
export const eqItemStatusList = (a: ItemStatus[], b: ItemStatus[]): boolean => rt.eqList(a, b, eqItemStatusValue);

export const getItemStatusListBody = (buf: rt.Buffer, state: number): [ItemStatus[], rt.Err] => {
    const [list, err] = rt.getDefaultList<ItemStatus>(
        buf,
        state,
        () => DefaultItemStatus(),
        (buf) => {
            const [value, err] = rt.getU8(buf);
            if (err !== null) return [DefaultItemStatus(), rt.errU(err)];
            const item = value as ItemStatus;
            if (!IsItemStatus(item)) return [DefaultItemStatus(), new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
    return [list, err];
};

export const setItemStatusListBody = (buf: rt.Buffer, state: number, v: ItemStatus[]): rt.Err => {
    return rt.setDefaultList<ItemStatus>(
        buf,
        state,
        v,
        (item) => IsDefaultItemStatus(item),
        (buf, item) => {
            if (!IsItemStatus(item)) return new Error(`非法枚举值: ${item}`);
            return rt.setU8(buf, item as any);
        },
    );
};

// 可否选号
export enum SimPickPhone {
    No = 0,
    Yes = 1,
    Active = 3,
    Abcc = 4,
}

export const DefaultSimPickPhone = (): SimPickPhone => SimPickPhone.No;
export const IsSimPickPhone = (v: SimPickPhone): boolean => {
    switch (v) {
    case SimPickPhone.No:
    case SimPickPhone.Yes:
    case SimPickPhone.Active:
    case SimPickPhone.Abcc:
        return true;
    default:
        return false;
    }
};

export const NormalizeSimPickPhone = (v: SimPickPhone): SimPickPhone => {
    if (IsSimPickPhone(v)) return v;
    if ((v as any) === 0) return DefaultSimPickPhone();
    return v;
};
export const IsDefaultSimPickPhone = (v: SimPickPhone): boolean => NormalizeSimPickPhone(v) === DefaultSimPickPhone();
export const IsAssignableSimPickPhone = (v: SimPickPhone): boolean => IsSimPickPhone(v) || (((v as any) === 0) && !IsSimPickPhone(0 as any));
export const eqSimPickPhoneValue = (a: SimPickPhone, b: SimPickPhone): boolean => NormalizeSimPickPhone(a) === NormalizeSimPickPhone(b);
export const eqSimPickPhoneList = (a: SimPickPhone[], b: SimPickPhone[]): boolean => rt.eqList(a, b, eqSimPickPhoneValue);

export const getSimPickPhoneListBody = (buf: rt.Buffer, state: number): [SimPickPhone[], rt.Err] => {
    const [list, err] = rt.getDefaultList<SimPickPhone>(
        buf,
        state,
        () => DefaultSimPickPhone(),
        (buf) => {
            const [value, err] = rt.getU8(buf);
            if (err !== null) return [DefaultSimPickPhone(), rt.errU(err)];
            const item = value as SimPickPhone;
            if (!IsSimPickPhone(item)) return [DefaultSimPickPhone(), new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
    return [list, err];
};

export const setSimPickPhoneListBody = (buf: rt.Buffer, state: number, v: SimPickPhone[]): rt.Err => {
    return rt.setDefaultList<SimPickPhone>(
        buf,
        state,
        v,
        (item) => IsDefaultSimPickPhone(item),
        (buf, item) => {
            if (!IsSimPickPhone(item)) return new Error(`非法枚举值: ${item}`);
            return rt.setU8(buf, item as any);
        },
    );
};

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

export const DefaultSimOperator = (): SimOperator => SimOperator.Zz;
export const IsSimOperator = (v: SimOperator): boolean => {
    switch (v) {
    case SimOperator.Zz:
    case SimOperator.Lt:
    case SimOperator.Yd:
    case SimOperator.Dx:
    case SimOperator.Gd:
    case SimOperator.Xx:
    case SimOperator.A:
    case SimOperator.B:
        return true;
    default:
        return false;
    }
};

export const NormalizeSimOperator = (v: SimOperator): SimOperator => {
    if (IsSimOperator(v)) return v;
    if ((v as any) === 0) return DefaultSimOperator();
    return v;
};
export const IsDefaultSimOperator = (v: SimOperator): boolean => NormalizeSimOperator(v) === DefaultSimOperator();
export const IsAssignableSimOperator = (v: SimOperator): boolean => IsSimOperator(v) || (((v as any) === 0) && !IsSimOperator(0 as any));
export const eqSimOperatorValue = (a: SimOperator, b: SimOperator): boolean => NormalizeSimOperator(a) === NormalizeSimOperator(b);
export const eqSimOperatorList = (a: SimOperator[], b: SimOperator[]): boolean => rt.eqList(a, b, eqSimOperatorValue);

export const getSimOperatorListBody = (buf: rt.Buffer, state: number): [SimOperator[], rt.Err] => {
    const [list, err] = rt.getDefaultList<SimOperator>(
        buf,
        state,
        () => DefaultSimOperator(),
        (buf) => {
            const [value, err] = rt.getU8(buf);
            if (err !== null) return [DefaultSimOperator(), rt.errU(err)];
            const item = value as SimOperator;
            if (!IsSimOperator(item)) return [DefaultSimOperator(), new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
    return [list, err];
};

export const setSimOperatorListBody = (buf: rt.Buffer, state: number, v: SimOperator[]): rt.Err => {
    return rt.setDefaultList<SimOperator>(
        buf,
        state,
        v,
        (item) => IsDefaultSimOperator(item),
        (buf, item) => {
            if (!IsSimOperator(item)) return new Error(`非法枚举值: ${item}`);
            return rt.setU8(buf, item as any);
        },
    );
};

// 订单状态
// Pending      待处理
// Closed       已关闭
// Canceled     已取消
// Shipped      已发货
// Delivered    已送达
// Actived      已激活
// Settled      已结算
export enum OrderStatus {
    Pending = 0,
    Closed = 1,
    Canceled = 2,
    Shipped = 3,
    Delivered = 4,
    Actived = 5,
    Settled = 6,
}

export const DefaultOrderStatus = (): OrderStatus => OrderStatus.Pending;
export const IsOrderStatus = (v: OrderStatus): boolean => {
    switch (v) {
    case OrderStatus.Pending:
    case OrderStatus.Closed:
    case OrderStatus.Canceled:
    case OrderStatus.Shipped:
    case OrderStatus.Delivered:
    case OrderStatus.Actived:
    case OrderStatus.Settled:
        return true;
    default:
        return false;
    }
};

export const NormalizeOrderStatus = (v: OrderStatus): OrderStatus => {
    if (IsOrderStatus(v)) return v;
    if ((v as any) === 0) return DefaultOrderStatus();
    return v;
};
export const IsDefaultOrderStatus = (v: OrderStatus): boolean => NormalizeOrderStatus(v) === DefaultOrderStatus();
export const IsAssignableOrderStatus = (v: OrderStatus): boolean => IsOrderStatus(v) || (((v as any) === 0) && !IsOrderStatus(0 as any));
export const eqOrderStatusValue = (a: OrderStatus, b: OrderStatus): boolean => NormalizeOrderStatus(a) === NormalizeOrderStatus(b);
export const eqOrderStatusList = (a: OrderStatus[], b: OrderStatus[]): boolean => rt.eqList(a, b, eqOrderStatusValue);

export const getOrderStatusListBody = (buf: rt.Buffer, state: number): [OrderStatus[], rt.Err] => {
    const [list, err] = rt.getDefaultList<OrderStatus>(
        buf,
        state,
        () => DefaultOrderStatus(),
        (buf) => {
            const [value, err] = rt.getU8(buf);
            if (err !== null) return [DefaultOrderStatus(), rt.errU(err)];
            const item = value as OrderStatus;
            if (!IsOrderStatus(item)) return [DefaultOrderStatus(), new Error(`非法枚举值: ${item}`)];
            return [item, undefined];
        },
    );
    return [list, err];
};

export const setOrderStatusListBody = (buf: rt.Buffer, state: number, v: OrderStatus[]): rt.Err => {
    return rt.setDefaultList<OrderStatus>(
        buf,
        state,
        v,
        (item) => IsDefaultOrderStatus(item),
        (buf, item) => {
            if (!IsOrderStatus(item)) return new Error(`非法枚举值: ${item}`);
            return rt.setU8(buf, item as any);
        },
    );
};

