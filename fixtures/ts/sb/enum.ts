import * as rt from "./runtime_core"
import * as rm from "./runtime_meta"

// 账户状态
export enum AccountStatus {
    Offline = 0,
    Online = 1,
    Deleted = 2,
}

const accountStatusMeta = rm.defineEnum<AccountStatus>(AccountStatus.Offline, [
    AccountStatus.Offline,
    AccountStatus.Online,
    AccountStatus.Deleted,
] as const);

export const DefaultAccountStatus = (): AccountStatus => accountStatusMeta.defaultValue;
export const IsAccountStatus = (v: AccountStatus): boolean => rm.isEnum(accountStatusMeta, v);
export const NormalizeAccountStatus = (v: AccountStatus): AccountStatus => rm.normalizeEnum(accountStatusMeta, v);
export const IsDefaultAccountStatus = (v: AccountStatus): boolean => rm.isDefaultEnum(accountStatusMeta, v);
export const IsAssignableAccountStatus = (v: AccountStatus): boolean => rm.isAssignableEnum(accountStatusMeta, v);
export const eqAccountStatusValue = (a: AccountStatus, b: AccountStatus): boolean => rm.eqEnumValue(accountStatusMeta, a, b);
export const eqAccountStatusList = (a: AccountStatus[], b: AccountStatus[]): boolean => rm.eqEnumList(accountStatusMeta, a, b);

export const getAccountStatusListBody = (buf: rt.Buffer, state: number): [AccountStatus[], rt.Err] => rm.getEnumList(accountStatusMeta, buf, state);

export const setAccountStatusListBody = (buf: rt.Buffer, state: number, v: AccountStatus[]): rt.Err => rm.setEnumList(accountStatusMeta, buf, state, v);

// 类型
export enum Type {
    Sim = 0,
    Recharge = 1,
}

const typeMeta = rm.defineEnum<Type>(Type.Sim, [
    Type.Sim,
    Type.Recharge,
] as const);

export const DefaultType = (): Type => typeMeta.defaultValue;
export const IsType = (v: Type): boolean => rm.isEnum(typeMeta, v);
export const NormalizeType = (v: Type): Type => rm.normalizeEnum(typeMeta, v);
export const IsDefaultType = (v: Type): boolean => rm.isDefaultEnum(typeMeta, v);
export const IsAssignableType = (v: Type): boolean => rm.isAssignableEnum(typeMeta, v);
export const eqTypeValue = (a: Type, b: Type): boolean => rm.eqEnumValue(typeMeta, a, b);
export const eqTypeList = (a: Type[], b: Type[]): boolean => rm.eqEnumList(typeMeta, a, b);

export const getTypeListBody = (buf: rt.Buffer, state: number): [Type[], rt.Err] => rm.getEnumList(typeMeta, buf, state);

export const setTypeListBody = (buf: rt.Buffer, state: number, v: Type[]): rt.Err => rm.setEnumList(typeMeta, buf, state, v);

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

const statusMeta = rm.defineEnum<Status>(Status.Ok, [
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
export const IsStatus = (v: Status): boolean => rm.isEnum(statusMeta, v);
export const NormalizeStatus = (v: Status): Status => rm.normalizeEnum(statusMeta, v);
export const IsDefaultStatus = (v: Status): boolean => rm.isDefaultEnum(statusMeta, v);
export const IsAssignableStatus = (v: Status): boolean => rm.isAssignableEnum(statusMeta, v);
export const eqStatusValue = (a: Status, b: Status): boolean => rm.eqEnumValue(statusMeta, a, b);
export const eqStatusList = (a: Status[], b: Status[]): boolean => rm.eqEnumList(statusMeta, a, b);

export const getStatusListBody = (buf: rt.Buffer, state: number): [Status[], rt.Err] => rm.getEnumList(statusMeta, buf, state);

export const setStatusListBody = (buf: rt.Buffer, state: number, v: Status[]): rt.Err => rm.setEnumList(statusMeta, buf, state, v);

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

const statusAMeta = rm.defineEnum<StatusA>(StatusA.Ok, [
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
export const IsStatusA = (v: StatusA): boolean => rm.isEnum(statusAMeta, v);
export const NormalizeStatusA = (v: StatusA): StatusA => rm.normalizeEnum(statusAMeta, v);
export const IsDefaultStatusA = (v: StatusA): boolean => rm.isDefaultEnum(statusAMeta, v);
export const IsAssignableStatusA = (v: StatusA): boolean => rm.isAssignableEnum(statusAMeta, v);
export const eqStatusAValue = (a: StatusA, b: StatusA): boolean => rm.eqEnumValue(statusAMeta, a, b);
export const eqStatusAList = (a: StatusA[], b: StatusA[]): boolean => rm.eqEnumList(statusAMeta, a, b);

export const getStatusAListBody = (buf: rt.Buffer, state: number): [StatusA[], rt.Err] => rm.getEnumList(statusAMeta, buf, state);

export const setStatusAListBody = (buf: rt.Buffer, state: number, v: StatusA[]): rt.Err => rm.setEnumList(statusAMeta, buf, state, v);

// 订单状态
export enum ItemStatus {
    Offline = 0,
    Online = 1,
}

const itemStatusMeta = rm.defineEnum<ItemStatus>(ItemStatus.Offline, [
    ItemStatus.Offline,
    ItemStatus.Online,
] as const);

export const DefaultItemStatus = (): ItemStatus => itemStatusMeta.defaultValue;
export const IsItemStatus = (v: ItemStatus): boolean => rm.isEnum(itemStatusMeta, v);
export const NormalizeItemStatus = (v: ItemStatus): ItemStatus => rm.normalizeEnum(itemStatusMeta, v);
export const IsDefaultItemStatus = (v: ItemStatus): boolean => rm.isDefaultEnum(itemStatusMeta, v);
export const IsAssignableItemStatus = (v: ItemStatus): boolean => rm.isAssignableEnum(itemStatusMeta, v);
export const eqItemStatusValue = (a: ItemStatus, b: ItemStatus): boolean => rm.eqEnumValue(itemStatusMeta, a, b);
export const eqItemStatusList = (a: ItemStatus[], b: ItemStatus[]): boolean => rm.eqEnumList(itemStatusMeta, a, b);

export const getItemStatusListBody = (buf: rt.Buffer, state: number): [ItemStatus[], rt.Err] => rm.getEnumList(itemStatusMeta, buf, state);

export const setItemStatusListBody = (buf: rt.Buffer, state: number, v: ItemStatus[]): rt.Err => rm.setEnumList(itemStatusMeta, buf, state, v);

// 可否选号
export enum SimPickPhone {
    No = 0,
    Yes = 1,
    Active = 3,
    Abcc = 4,
}

const simPickPhoneMeta = rm.defineEnum<SimPickPhone>(SimPickPhone.No, [
    SimPickPhone.No,
    SimPickPhone.Yes,
    SimPickPhone.Active,
    SimPickPhone.Abcc,
] as const);

export const DefaultSimPickPhone = (): SimPickPhone => simPickPhoneMeta.defaultValue;
export const IsSimPickPhone = (v: SimPickPhone): boolean => rm.isEnum(simPickPhoneMeta, v);
export const NormalizeSimPickPhone = (v: SimPickPhone): SimPickPhone => rm.normalizeEnum(simPickPhoneMeta, v);
export const IsDefaultSimPickPhone = (v: SimPickPhone): boolean => rm.isDefaultEnum(simPickPhoneMeta, v);
export const IsAssignableSimPickPhone = (v: SimPickPhone): boolean => rm.isAssignableEnum(simPickPhoneMeta, v);
export const eqSimPickPhoneValue = (a: SimPickPhone, b: SimPickPhone): boolean => rm.eqEnumValue(simPickPhoneMeta, a, b);
export const eqSimPickPhoneList = (a: SimPickPhone[], b: SimPickPhone[]): boolean => rm.eqEnumList(simPickPhoneMeta, a, b);

export const getSimPickPhoneListBody = (buf: rt.Buffer, state: number): [SimPickPhone[], rt.Err] => rm.getEnumList(simPickPhoneMeta, buf, state);

export const setSimPickPhoneListBody = (buf: rt.Buffer, state: number, v: SimPickPhone[]): rt.Err => rm.setEnumList(simPickPhoneMeta, buf, state, v);

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

const simOperatorMeta = rm.defineEnum<SimOperator>(SimOperator.Zz, [
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
export const IsSimOperator = (v: SimOperator): boolean => rm.isEnum(simOperatorMeta, v);
export const NormalizeSimOperator = (v: SimOperator): SimOperator => rm.normalizeEnum(simOperatorMeta, v);
export const IsDefaultSimOperator = (v: SimOperator): boolean => rm.isDefaultEnum(simOperatorMeta, v);
export const IsAssignableSimOperator = (v: SimOperator): boolean => rm.isAssignableEnum(simOperatorMeta, v);
export const eqSimOperatorValue = (a: SimOperator, b: SimOperator): boolean => rm.eqEnumValue(simOperatorMeta, a, b);
export const eqSimOperatorList = (a: SimOperator[], b: SimOperator[]): boolean => rm.eqEnumList(simOperatorMeta, a, b);

export const getSimOperatorListBody = (buf: rt.Buffer, state: number): [SimOperator[], rt.Err] => rm.getEnumList(simOperatorMeta, buf, state);

export const setSimOperatorListBody = (buf: rt.Buffer, state: number, v: SimOperator[]): rt.Err => rm.setEnumList(simOperatorMeta, buf, state, v);

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

const orderStatusMeta = rm.defineEnum<OrderStatus>(OrderStatus.Pending, [
    OrderStatus.Pending,
    OrderStatus.Closed,
    OrderStatus.Canceled,
    OrderStatus.Shipped,
    OrderStatus.Delivered,
    OrderStatus.Actived,
    OrderStatus.Settled,
] as const);

export const DefaultOrderStatus = (): OrderStatus => orderStatusMeta.defaultValue;
export const IsOrderStatus = (v: OrderStatus): boolean => rm.isEnum(orderStatusMeta, v);
export const NormalizeOrderStatus = (v: OrderStatus): OrderStatus => rm.normalizeEnum(orderStatusMeta, v);
export const IsDefaultOrderStatus = (v: OrderStatus): boolean => rm.isDefaultEnum(orderStatusMeta, v);
export const IsAssignableOrderStatus = (v: OrderStatus): boolean => rm.isAssignableEnum(orderStatusMeta, v);
export const eqOrderStatusValue = (a: OrderStatus, b: OrderStatus): boolean => rm.eqEnumValue(orderStatusMeta, a, b);
export const eqOrderStatusList = (a: OrderStatus[], b: OrderStatus[]): boolean => rm.eqEnumList(orderStatusMeta, a, b);

export const getOrderStatusListBody = (buf: rt.Buffer, state: number): [OrderStatus[], rt.Err] => rm.getEnumList(orderStatusMeta, buf, state);

export const setOrderStatusListBody = (buf: rt.Buffer, state: number, v: OrderStatus[]): rt.Err => rm.setEnumList(orderStatusMeta, buf, state, v);

