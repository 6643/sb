import * as rt from "./type"

// 账户状态
export enum AccountStatus {
    Offline = 0,
    Online = 1,
    Deleted = 2,
}

const accountStatusMeta = rt.defineEnum<AccountStatus>(AccountStatus.Offline, [
    AccountStatus.Offline,
    AccountStatus.Online,
    AccountStatus.Deleted,
] as const);

export const DefaultAccountStatus = (): AccountStatus => accountStatusMeta.defaultValue;
export const IsAccountStatus = (v: AccountStatus): boolean => rt.isEnum(accountStatusMeta, v);
export const NormalizeAccountStatus = (v: AccountStatus): AccountStatus => rt.normalizeEnum(accountStatusMeta, v);
export const IsDefaultAccountStatus = (v: AccountStatus): boolean => rt.isDefaultEnum(accountStatusMeta, v);
export const IsAssignableAccountStatus = (v: AccountStatus): boolean => rt.isAssignableEnum(accountStatusMeta, v);
export const eqAccountStatusValue = (a: AccountStatus, b: AccountStatus): boolean => rt.eqEnumValue(accountStatusMeta, a, b);
export const eqAccountStatusList = (a: AccountStatus[], b: AccountStatus[]): boolean => rt.eqEnumList(accountStatusMeta, a, b);

export const getAccountStatusListBody = (buf: rt.Buffer, state: number): [AccountStatus[], rt.Err] => rt.getEnumList(accountStatusMeta, buf, state);

export const setAccountStatusListBody = (buf: rt.Buffer, state: number, v: AccountStatus[]): rt.Err => rt.setEnumList(accountStatusMeta, buf, state, v);

// 类型
export enum Type {
    Sim = 0,
    Recharge = 1,
}

const typeMeta = rt.defineEnum<Type>(Type.Sim, [
    Type.Sim,
    Type.Recharge,
] as const);

export const DefaultType = (): Type => typeMeta.defaultValue;
export const IsType = (v: Type): boolean => rt.isEnum(typeMeta, v);
export const NormalizeType = (v: Type): Type => rt.normalizeEnum(typeMeta, v);
export const IsDefaultType = (v: Type): boolean => rt.isDefaultEnum(typeMeta, v);
export const IsAssignableType = (v: Type): boolean => rt.isAssignableEnum(typeMeta, v);
export const eqTypeValue = (a: Type, b: Type): boolean => rt.eqEnumValue(typeMeta, a, b);
export const eqTypeList = (a: Type[], b: Type[]): boolean => rt.eqEnumList(typeMeta, a, b);

export const getTypeListBody = (buf: rt.Buffer, state: number): [Type[], rt.Err] => rt.getEnumList(typeMeta, buf, state);

export const setTypeListBody = (buf: rt.Buffer, state: number, v: Type[]): rt.Err => rt.setEnumList(typeMeta, buf, state, v);

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

const statusMeta = rt.defineEnum<Status>(Status.Ok, [
    Status.Ok,
    Status.Err,
    Status.Two,
    Status.Three,
    Status.Four,
    Status.Five,
    Status.Six,
    Status.Seven,
    Status.One,
] as const);

export const DefaultStatus = (): Status => statusMeta.defaultValue;
export const IsStatus = (v: Status): boolean => rt.isEnum(statusMeta, v);
export const NormalizeStatus = (v: Status): Status => rt.normalizeEnum(statusMeta, v);
export const IsDefaultStatus = (v: Status): boolean => rt.isDefaultEnum(statusMeta, v);
export const IsAssignableStatus = (v: Status): boolean => rt.isAssignableEnum(statusMeta, v);
export const eqStatusValue = (a: Status, b: Status): boolean => rt.eqEnumValue(statusMeta, a, b);
export const eqStatusList = (a: Status[], b: Status[]): boolean => rt.eqEnumList(statusMeta, a, b);

export const getStatusListBody = (buf: rt.Buffer, state: number): [Status[], rt.Err] => rt.getEnumList(statusMeta, buf, state);

export const setStatusListBody = (buf: rt.Buffer, state: number, v: Status[]): rt.Err => rt.setEnumList(statusMeta, buf, state, v);

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

const statusAMeta = rt.defineEnum<StatusA>(StatusA.Ok, [
    StatusA.Ok,
    StatusA.One,
    StatusA.Two,
    StatusA.Three,
    StatusA.Four,
    StatusA.Five,
    StatusA.Six,
    StatusA.Seven,
] as const);

export const DefaultStatusA = (): StatusA => statusAMeta.defaultValue;
export const IsStatusA = (v: StatusA): boolean => rt.isEnum(statusAMeta, v);
export const NormalizeStatusA = (v: StatusA): StatusA => rt.normalizeEnum(statusAMeta, v);
export const IsDefaultStatusA = (v: StatusA): boolean => rt.isDefaultEnum(statusAMeta, v);
export const IsAssignableStatusA = (v: StatusA): boolean => rt.isAssignableEnum(statusAMeta, v);
export const eqStatusAValue = (a: StatusA, b: StatusA): boolean => rt.eqEnumValue(statusAMeta, a, b);
export const eqStatusAList = (a: StatusA[], b: StatusA[]): boolean => rt.eqEnumList(statusAMeta, a, b);

export const getStatusAListBody = (buf: rt.Buffer, state: number): [StatusA[], rt.Err] => rt.getEnumList(statusAMeta, buf, state);

export const setStatusAListBody = (buf: rt.Buffer, state: number, v: StatusA[]): rt.Err => rt.setEnumList(statusAMeta, buf, state, v);

// 订单状态
export enum ItemStatus {
    Offline = 0,
    Online = 1,
}

const itemStatusMeta = rt.defineEnum<ItemStatus>(ItemStatus.Offline, [
    ItemStatus.Offline,
    ItemStatus.Online,
] as const);

export const DefaultItemStatus = (): ItemStatus => itemStatusMeta.defaultValue;
export const IsItemStatus = (v: ItemStatus): boolean => rt.isEnum(itemStatusMeta, v);
export const NormalizeItemStatus = (v: ItemStatus): ItemStatus => rt.normalizeEnum(itemStatusMeta, v);
export const IsDefaultItemStatus = (v: ItemStatus): boolean => rt.isDefaultEnum(itemStatusMeta, v);
export const IsAssignableItemStatus = (v: ItemStatus): boolean => rt.isAssignableEnum(itemStatusMeta, v);
export const eqItemStatusValue = (a: ItemStatus, b: ItemStatus): boolean => rt.eqEnumValue(itemStatusMeta, a, b);
export const eqItemStatusList = (a: ItemStatus[], b: ItemStatus[]): boolean => rt.eqEnumList(itemStatusMeta, a, b);

export const getItemStatusListBody = (buf: rt.Buffer, state: number): [ItemStatus[], rt.Err] => rt.getEnumList(itemStatusMeta, buf, state);

export const setItemStatusListBody = (buf: rt.Buffer, state: number, v: ItemStatus[]): rt.Err => rt.setEnumList(itemStatusMeta, buf, state, v);

// 可否选号
export enum SimPickPhone {
    No = 0,
    Yes = 1,
    Active = 3,
    Abcc = 4,
}

const simPickPhoneMeta = rt.defineEnum<SimPickPhone>(SimPickPhone.No, [
    SimPickPhone.No,
    SimPickPhone.Yes,
    SimPickPhone.Active,
    SimPickPhone.Abcc,
] as const);

export const DefaultSimPickPhone = (): SimPickPhone => simPickPhoneMeta.defaultValue;
export const IsSimPickPhone = (v: SimPickPhone): boolean => rt.isEnum(simPickPhoneMeta, v);
export const NormalizeSimPickPhone = (v: SimPickPhone): SimPickPhone => rt.normalizeEnum(simPickPhoneMeta, v);
export const IsDefaultSimPickPhone = (v: SimPickPhone): boolean => rt.isDefaultEnum(simPickPhoneMeta, v);
export const IsAssignableSimPickPhone = (v: SimPickPhone): boolean => rt.isAssignableEnum(simPickPhoneMeta, v);
export const eqSimPickPhoneValue = (a: SimPickPhone, b: SimPickPhone): boolean => rt.eqEnumValue(simPickPhoneMeta, a, b);
export const eqSimPickPhoneList = (a: SimPickPhone[], b: SimPickPhone[]): boolean => rt.eqEnumList(simPickPhoneMeta, a, b);

export const getSimPickPhoneListBody = (buf: rt.Buffer, state: number): [SimPickPhone[], rt.Err] => rt.getEnumList(simPickPhoneMeta, buf, state);

export const setSimPickPhoneListBody = (buf: rt.Buffer, state: number, v: SimPickPhone[]): rt.Err => rt.setEnumList(simPickPhoneMeta, buf, state, v);

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

const simOperatorMeta = rt.defineEnum<SimOperator>(SimOperator.Zz, [
    SimOperator.Zz,
    SimOperator.Lt,
    SimOperator.Yd,
    SimOperator.Dx,
    SimOperator.Gd,
    SimOperator.Xx,
    SimOperator.A,
    SimOperator.B,
] as const);

export const DefaultSimOperator = (): SimOperator => simOperatorMeta.defaultValue;
export const IsSimOperator = (v: SimOperator): boolean => rt.isEnum(simOperatorMeta, v);
export const NormalizeSimOperator = (v: SimOperator): SimOperator => rt.normalizeEnum(simOperatorMeta, v);
export const IsDefaultSimOperator = (v: SimOperator): boolean => rt.isDefaultEnum(simOperatorMeta, v);
export const IsAssignableSimOperator = (v: SimOperator): boolean => rt.isAssignableEnum(simOperatorMeta, v);
export const eqSimOperatorValue = (a: SimOperator, b: SimOperator): boolean => rt.eqEnumValue(simOperatorMeta, a, b);
export const eqSimOperatorList = (a: SimOperator[], b: SimOperator[]): boolean => rt.eqEnumList(simOperatorMeta, a, b);

export const getSimOperatorListBody = (buf: rt.Buffer, state: number): [SimOperator[], rt.Err] => rt.getEnumList(simOperatorMeta, buf, state);

export const setSimOperatorListBody = (buf: rt.Buffer, state: number, v: SimOperator[]): rt.Err => rt.setEnumList(simOperatorMeta, buf, state, v);

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

const orderStatusMeta = rt.defineEnum<OrderStatus>(OrderStatus.Pending, [
    OrderStatus.Pending,
    OrderStatus.Closed,
    OrderStatus.Canceled,
    OrderStatus.Shipped,
    OrderStatus.Delivered,
    OrderStatus.Actived,
    OrderStatus.Settled,
] as const);

export const DefaultOrderStatus = (): OrderStatus => orderStatusMeta.defaultValue;
export const IsOrderStatus = (v: OrderStatus): boolean => rt.isEnum(orderStatusMeta, v);
export const NormalizeOrderStatus = (v: OrderStatus): OrderStatus => rt.normalizeEnum(orderStatusMeta, v);
export const IsDefaultOrderStatus = (v: OrderStatus): boolean => rt.isDefaultEnum(orderStatusMeta, v);
export const IsAssignableOrderStatus = (v: OrderStatus): boolean => rt.isAssignableEnum(orderStatusMeta, v);
export const eqOrderStatusValue = (a: OrderStatus, b: OrderStatus): boolean => rt.eqEnumValue(orderStatusMeta, a, b);
export const eqOrderStatusList = (a: OrderStatus[], b: OrderStatus[]): boolean => rt.eqEnumList(orderStatusMeta, a, b);

export const getOrderStatusListBody = (buf: rt.Buffer, state: number): [OrderStatus[], rt.Err] => rt.getEnumList(orderStatusMeta, buf, state);

export const setOrderStatusListBody = (buf: rt.Buffer, state: number, v: OrderStatus[]): rt.Err => rt.setEnumList(orderStatusMeta, buf, state, v);

