import * as rt from "./type"
import { DefaultOrderStatus, IsAssignableOrderStatus, IsDefaultOrderStatus, IsOrderStatus, NormalizeOrderStatus, OrderStatus, getOrderStatusListBody, setOrderStatusListBody } from "./enum"
import { SimInfo, eqSimInfo, getSimInfoListBody, isZeroSimInfo, newSimInfo, readSimInfo, setSimInfo, setSimInfoListBody, validateSimInfo } from "./struct_sim_info"

export interface RechargeB extends rt.Serializable, rt.Deserializable {
    // abcd
    id: number;
    type: OrderStatus[];
    phone: string[];
    si: SimInfo | null;
    bid: number;
}

const rechargeBHeaderWidths = [1, 2, 2, 1, 1] as const;

const rechargeBMeta = rt.defineStruct<RechargeB>({
    name: "RechargeB",
    headerWidths: rechargeBHeaderWidths,
    create: () => ({
        id: 0,
        type: [],
        phone: [],
        si: null,
        bid: 0,
    }) as any as RechargeB,
    fields: [

        rt.scalarField<RechargeB, number>("id", "Id", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.defaultListField<RechargeB, OrderStatus>("type", "Type", DefaultOrderStatus, getOrderStatusListBody, setOrderStatusListBody, (item) => IsAssignableOrderStatus(item as any) ? undefined : new Error(`非法枚举值: ${item as any}`), (item) => IsDefaultOrderStatus(item as any), (left, right) => NormalizeOrderStatus(left as any) === NormalizeOrderStatus(right as any)),
        rt.textListField<RechargeB>("phone", "Phone"),
        rt.structField<RechargeB, SimInfo>("si", "Si", isZeroSimInfo, readSimInfo, setSimInfo, validateSimInfo, eqSimInfo),
        rt.scalarField<RechargeB, number>("bid", "Bid", 0, rt.getU32, rt.setU32, rt.eqU32),
    ],
});

export const newRechargeB = (): RechargeB => rt.newStruct(rechargeBMeta, getRechargeB, setRechargeB);
export const isZeroRechargeB = (s: RechargeB | null | undefined): boolean => rt.isZeroStruct(rechargeBMeta, s);
export const validateRechargeB = (s: RechargeB | null | undefined): rt.Err => rt.validateStruct(rechargeBMeta, s);
export const getRechargeB = (buf: rt.Buffer): [RechargeB, rt.Err] => rt.getStruct(rechargeBMeta, buf);
export const setRechargeB = (buf: rt.Buffer, s: RechargeB): rt.Err => rt.setStruct(rechargeBMeta, buf, s);
export const readRechargeB = (buf: rt.Buffer): [RechargeB, rt.Err] => getRechargeB(buf);
export const eqRechargeB = (a: RechargeB | null | undefined, b: RechargeB | null | undefined): boolean => rt.eqStruct(rechargeBMeta, a, b);

export const getRechargeBListBody = (buf: rt.Buffer, state: number): [RechargeB[], rt.Err] => {
    const [list, err] = rt.getDefaultList<RechargeB>(
        buf,
        state,
        () => newRechargeB(),
        (buf) => readRechargeB(buf),
    );
    return [list, err];
};

export const setRechargeBListBody = (buf: rt.Buffer, state: number, v: RechargeB[]): rt.Err => {
    return rt.setDefaultList<RechargeB>(
        buf,
        state,
        v,
        (item) => isZeroRechargeB(item),
        (buf, item) => setRechargeB(buf, item),
    );
}
