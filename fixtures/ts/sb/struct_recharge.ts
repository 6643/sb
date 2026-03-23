import * as rt from "./runtime_core"
import * as rm from "./runtime_meta"
import { DefaultOrderStatus, IsAssignableOrderStatus, IsDefaultOrderStatus, IsOrderStatus, NormalizeOrderStatus, OrderStatus, getOrderStatusListBody, setOrderStatusListBody } from "./enum"
import { SimInfo, eqSimInfo, getSimInfoListBody, isZeroSimInfo, newSimInfo, readSimInfo, setSimInfo, setSimInfoListBody, validateSimInfo } from "./struct_sim_info"

export interface Recharge extends rt.Serializable, rt.Deserializable {
    // abcd
    id: number;
    type: OrderStatus[];
    phone: string[];
    si: SimInfo | null;
}

const rechargeHeaderWidths = [1, 2, 2, 1] as const;

const rechargeMeta = rm.defineStruct<Recharge>({
    name: "Recharge",
    headerWidths: rechargeHeaderWidths,
    create: () => ({
        id: 0,
        type: [],
        phone: [],
        si: null,
    }) as any as Recharge,
    fields: [

        rm.scalarField<Recharge, number>("id", "Id", 0, rt.getU32, rt.setU32, rt.eqU32),
        rm.defaultListField<Recharge, OrderStatus>("type", "Type", DefaultOrderStatus, getOrderStatusListBody, setOrderStatusListBody, (item) => IsAssignableOrderStatus(item as any) ? undefined : new Error(`非法枚举值: ${item as any}`), (item) => IsDefaultOrderStatus(item as any), (left, right) => NormalizeOrderStatus(left as any) === NormalizeOrderStatus(right as any)),
        rm.textListField<Recharge>("phone", "Phone"),
        rm.structField<Recharge, SimInfo>("si", "Si", isZeroSimInfo, readSimInfo, setSimInfo, validateSimInfo, eqSimInfo),
    ],
});

export const newRecharge = (): Recharge => rm.newStruct(rechargeMeta, getRecharge, setRecharge);
export const isZeroRecharge = (s: Recharge | null | undefined): boolean => rm.isZeroStruct(rechargeMeta, s);
export const validateRecharge = (s: Recharge | null | undefined): rt.Err => rm.validateStruct(rechargeMeta, s);
export const getRecharge = (buf: rt.Buffer): [Recharge, rt.Err] => rm.getStruct(rechargeMeta, buf);
export const setRecharge = (buf: rt.Buffer, s: Recharge): rt.Err => rm.setStruct(rechargeMeta, buf, s);
export const readRecharge = (buf: rt.Buffer): [Recharge, rt.Err] => getRecharge(buf);
export const eqRecharge = (a: Recharge | null | undefined, b: Recharge | null | undefined): boolean => rm.eqStruct(rechargeMeta, a, b);

export const getRechargeListBody = (buf: rt.Buffer, state: number): [Recharge[], rt.Err] => {
    const [list, err] = rt.getDefaultList<Recharge>(
        buf,
        state,
        () => newRecharge(),
        (buf) => readRecharge(buf),
    );
    return [list, err];
};

export const setRechargeListBody = (buf: rt.Buffer, state: number, v: Recharge[]): rt.Err => {
    return rt.setDefaultList<Recharge>(
        buf,
        state,
        v,
        (item) => isZeroRecharge(item),
        (buf, item) => setRecharge(buf, item),
    );
}
