import * as rt from "./type"
import { DefaultOrderStatus, IsAssignableOrderStatus, IsDefaultOrderStatus, IsOrderStatus, NormalizeOrderStatus, OrderStatus, getOrderStatusListBody, setOrderStatusListBody } from "./enum"

export interface SimOrder extends rt.Serializable, rt.Deserializable {
    id: number;
    accountId: number;
    itemId: number;
    // 办理人姓名
    name: string;
    // 联系电话
    phone: string;
    // 身份证号
    idNo: string;
    // 所在城市
    cityCode: number;
    // 详细地址
    address: string;
    // 新手机号码
    newPhone: string;
    // 佣金
    commission: number;
    status: OrderStatus;
}

const simOrderHeaderWidths = [1, 1, 1, 2, 2, 2, 1, 2, 2, 1, 1] as const;

const simOrderMeta = rt.defineStruct<SimOrder>({
    name: "SimOrder",
    headerWidths: simOrderHeaderWidths,
    create: () => ({
        id: 0,
        accountId: 0,
        itemId: 0,
        name: "",
        phone: "",
        idNo: "",
        cityCode: 0,
        address: "",
        newPhone: "",
        commission: 0,
        status: DefaultOrderStatus(),
    }) as any as SimOrder,
    fields: [

        rt.scalarField<SimOrder, number>("id", "Id", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.scalarField<SimOrder, number>("accountId", "AccountId", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.scalarField<SimOrder, number>("itemId", "ItemId", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.textField<SimOrder>("name", "Name"),
        rt.textField<SimOrder>("phone", "Phone"),
        rt.textField<SimOrder>("idNo", "IdNo"),
        rt.scalarField<SimOrder, number>("cityCode", "CityCode", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.textField<SimOrder>("address", "Address"),
        rt.textField<SimOrder>("newPhone", "NewPhone"),
        rt.scalarField<SimOrder, number>("commission", "Commission", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.enumField<SimOrder, OrderStatus>("status", "Status", DefaultOrderStatus, IsOrderStatus, IsAssignableOrderStatus, NormalizeOrderStatus),
    ],
});

export const newSimOrder = (): SimOrder => rt.newStruct(simOrderMeta, getSimOrder, setSimOrder);
export const isZeroSimOrder = (s: SimOrder | null | undefined): boolean => rt.isZeroStruct(simOrderMeta, s);
export const validateSimOrder = (s: SimOrder | null | undefined): rt.Err => rt.validateStruct(simOrderMeta, s);
export const getSimOrder = (buf: rt.Buffer): [SimOrder, rt.Err] => rt.getStruct(simOrderMeta, buf);
export const setSimOrder = (buf: rt.Buffer, s: SimOrder): rt.Err => rt.setStruct(simOrderMeta, buf, s);
export const readSimOrder = (buf: rt.Buffer): [SimOrder, rt.Err] => getSimOrder(buf);
export const eqSimOrder = (a: SimOrder | null | undefined, b: SimOrder | null | undefined): boolean => rt.eqStruct(simOrderMeta, a, b);

export const getSimOrderListBody = (buf: rt.Buffer, state: number): [SimOrder[], rt.Err] => {
    const [list, err] = rt.getDefaultList<SimOrder>(
        buf,
        state,
        () => newSimOrder(),
        (buf) => readSimOrder(buf),
    );
    return [list, err];
};

export const setSimOrderListBody = (buf: rt.Buffer, state: number, v: SimOrder[]): rt.Err => {
    return rt.setDefaultList<SimOrder>(
        buf,
        state,
        v,
        (item) => isZeroSimOrder(item),
        (buf, item) => setSimOrder(buf, item),
    );
}
