import * as rt from "./type"
import { DefaultOrderStatus, IsAssignableOrderStatus, IsDefaultOrderStatus, IsOrderStatus, NormalizeOrderStatus, OrderStatus, getOrderStatusListBody, setOrderStatusListBody } from "./enum"
import { SimInfo, eqSimInfo, getSimInfoListBody, isZeroSimInfo, newSimInfo, readSimInfo, setSimInfo, setSimInfoListBody, validateSimInfo } from "./struct_sim_info"

export interface RechargeA extends rt.Serializable, rt.Deserializable {
    // abcd
    id: number;
    type: OrderStatus[];
    phone: string[];
    si: SimInfo | null;
    aid: number;
}

const rechargeAHeaderWidths = [1, 2, 2, 1, 1] as const;

const rechargeAMeta = rt.defineStruct<RechargeA>({
    name: "RechargeA",
    headerWidths: rechargeAHeaderWidths,
    create: () => ({
        id: 0,
        type: [],
        phone: [],
        si: null,
        aid: 0,
    }) as any as RechargeA,
    fields: [

        rt.scalarField<RechargeA, number>("id", "Id", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.defaultListField<RechargeA, OrderStatus>("type", "Type", DefaultOrderStatus, getOrderStatusListBody, setOrderStatusListBody, (item) => IsAssignableOrderStatus(item as any) ? undefined : new Error(`非法枚举值: ${item as any}`), (item) => IsDefaultOrderStatus(item as any), (left, right) => NormalizeOrderStatus(left as any) === NormalizeOrderStatus(right as any)),
        rt.textListField<RechargeA>("phone", "Phone"),
        rt.structField<RechargeA, SimInfo>("si", "Si", isZeroSimInfo, readSimInfo, setSimInfo, validateSimInfo, eqSimInfo),
        rt.scalarField<RechargeA, number>("aid", "Aid", 0, rt.getU32, rt.setU32, rt.eqU32),
    ],
});

export const newRechargeA = (): RechargeA => rt.newStruct(rechargeAMeta, getRechargeA, setRechargeA);
export const isZeroRechargeA = (s: RechargeA | null | undefined): boolean => rt.isZeroStruct(rechargeAMeta, s);
export const validateRechargeA = (s: RechargeA | null | undefined): rt.Err => rt.validateStruct(rechargeAMeta, s);
export const getRechargeA = (buf: rt.Buffer): [RechargeA, rt.Err] => rt.getStruct(rechargeAMeta, buf);
export const setRechargeA = (buf: rt.Buffer, s: RechargeA): rt.Err => rt.setStruct(rechargeAMeta, buf, s);
export const readRechargeA = (buf: rt.Buffer): [RechargeA, rt.Err] => getRechargeA(buf);
export const eqRechargeA = (a: RechargeA | null | undefined, b: RechargeA | null | undefined): boolean => rt.eqStruct(rechargeAMeta, a, b);

export const getRechargeAListBody = (buf: rt.Buffer, state: number): [RechargeA[], rt.Err] => {
    const [list, err] = rt.getDefaultList<RechargeA>(
        buf,
        state,
        () => newRechargeA(),
        (buf) => readRechargeA(buf),
    );
    return [list, err];
};

export const setRechargeAListBody = (buf: rt.Buffer, state: number, v: RechargeA[]): rt.Err => {
    return rt.setDefaultList<RechargeA>(
        buf,
        state,
        v,
        (item) => isZeroRechargeA(item),
        (buf, item) => setRechargeA(buf, item),
    );
}
