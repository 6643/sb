import * as rt from "./type"
import { DefaultItemStatus, DefaultSimOperator, DefaultSimPickPhone, DefaultType, IsAssignableItemStatus, IsAssignableSimOperator, IsAssignableSimPickPhone, IsAssignableType, IsDefaultItemStatus, IsDefaultSimOperator, IsDefaultSimPickPhone, IsDefaultType, IsItemStatus, IsSimOperator, IsSimPickPhone, IsType, ItemStatus, NormalizeItemStatus, NormalizeSimOperator, NormalizeSimPickPhone, NormalizeType, SimOperator, SimPickPhone, Type, getItemStatusListBody, getSimOperatorListBody, getSimPickPhoneListBody, getTypeListBody, setItemStatusListBody, setSimOperatorListBody, setSimPickPhoneListBody, setTypeListBody } from "./enum"
import { SimInfo, eqSimInfo, getSimInfoListBody, isZeroSimInfo, newSimInfo, readSimInfo, setSimInfo, setSimInfoListBody, validateSimInfo } from "./struct_sim_info"

export interface Sim extends rt.Serializable, rt.Deserializable {
    // SIM卡ID
    id: number;
    type: Type;
    status: ItemStatus;
    // 佣金
    commission: number;
    // 供应商ID
    supplier: number;
    // 推广员ID
    aff: number;
    // 合约期(月), 0:长期
    contractDuration: number;
    name: string;
    // 运营商
    operator: SimOperator;
    // 月租
    monthly: number;
    // 通用流量
    flowUniversal: number;
    // 定向流量
    flowDirectional: number;
    // 流量是否结转
    canMoveFlow: boolean;
    // 每月通话(分钟)
    callMonth: number;
    callPrice: number;
    // 每月短信(条)
    smsMonth: number;
    smsPrice: number;
    minAge: number;
    maxAge: number;
    // 归属地, 0:随机, 1:收货地
    attribution: number;
    // 选号
    pickPhone: SimPickPhone[];
    // 首充渠道
    firstChargeLink: string;
    // 首充金额
    firstChargeMoney: string;
    // 首充返额
    firstChargeReturn: string;
    // 禁发区域
    banCity: number[];
    info: SimInfo[];
    // 套餐截图
    snapshot: string[];
}

const simHeaderWidths = [1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2] as const;

const simMeta = rt.defineStruct<Sim>({
    name: "Sim",
    headerWidths: simHeaderWidths,
    create: () => ({
        id: 0,
        type: DefaultType(),
        status: DefaultItemStatus(),
        commission: 0,
        supplier: 0,
        aff: 0,
        contractDuration: 0,
        name: "",
        operator: DefaultSimOperator(),
        monthly: 0,
        flowUniversal: 0,
        flowDirectional: 0,
        canMoveFlow: false,
        callMonth: 0,
        callPrice: 0,
        smsMonth: 0,
        smsPrice: 0,
        minAge: 0,
        maxAge: 0,
        attribution: 0,
        pickPhone: [],
        firstChargeLink: "",
        firstChargeMoney: "",
        firstChargeReturn: "",
        banCity: [],
        info: [],
        snapshot: [],
    }) as any as Sim,
    fields: [

        rt.scalarField<Sim, number>("id", "Id", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.enumField<Sim, Type>("type", "Type", DefaultType, IsType, IsAssignableType, NormalizeType),
        rt.enumField<Sim, ItemStatus>("status", "Status", DefaultItemStatus, IsItemStatus, IsAssignableItemStatus, NormalizeItemStatus),
        rt.scalarField<Sim, number>("commission", "Commission", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.scalarField<Sim, number>("supplier", "Supplier", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.scalarField<Sim, number>("aff", "Aff", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.scalarField<Sim, number>("contractDuration", "ContractDuration", 0, rt.getU8, rt.setU8, rt.eqU8),
        rt.textField<Sim>("name", "Name"),
        rt.enumField<Sim, SimOperator>("operator", "Operator", DefaultSimOperator, IsSimOperator, IsAssignableSimOperator, NormalizeSimOperator),
        rt.scalarField<Sim, number>("monthly", "Monthly", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.scalarField<Sim, number>("flowUniversal", "FlowUniversal", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.scalarField<Sim, number>("flowDirectional", "FlowDirectional", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.boolField<Sim>("canMoveFlow", "CanMoveFlow"),
        rt.scalarField<Sim, number>("callMonth", "CallMonth", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.scalarField<Sim, number>("callPrice", "CallPrice", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.scalarField<Sim, number>("smsMonth", "SmsMonth", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.scalarField<Sim, number>("smsPrice", "SmsPrice", 0, rt.getU16, rt.setU16, rt.eqU16),
        rt.scalarField<Sim, number>("minAge", "MinAge", 0, rt.getU8, rt.setU8, rt.eqU8),
        rt.scalarField<Sim, number>("maxAge", "MaxAge", 0, rt.getU8, rt.setU8, rt.eqU8),
        rt.scalarField<Sim, number>("attribution", "Attribution", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.defaultListField<Sim, SimPickPhone>("pickPhone", "PickPhone", DefaultSimPickPhone, getSimPickPhoneListBody, setSimPickPhoneListBody, (item) => IsAssignableSimPickPhone(item as any) ? undefined : new Error(`非法枚举值: ${item as any}`), (item) => IsDefaultSimPickPhone(item as any), (left, right) => NormalizeSimPickPhone(left as any) === NormalizeSimPickPhone(right as any)),
        rt.textField<Sim>("firstChargeLink", "FirstChargeLink"),
        rt.textField<Sim>("firstChargeMoney", "FirstChargeMoney"),
        rt.textField<Sim>("firstChargeReturn", "FirstChargeReturn"),
        rt.zeroListField<Sim, number>("banCity", "BanCity", 0, rt.getU32, rt.setU32, rt.eqU32),
        rt.defaultListField<Sim, SimInfo>("info", "Info", newSimInfo, getSimInfoListBody, setSimInfoListBody, validateSimInfo, isZeroSimInfo, eqSimInfo),
        rt.textListField<Sim>("snapshot", "Snapshot"),
    ],
});

export const newSim = (): Sim => rt.newStruct(simMeta, getSim, setSim);
export const isZeroSim = (s: Sim | null | undefined): boolean => rt.isZeroStruct(simMeta, s);
export const validateSim = (s: Sim | null | undefined): rt.Err => rt.validateStruct(simMeta, s);
export const getSim = (buf: rt.Buffer): [Sim, rt.Err] => rt.getStruct(simMeta, buf);
export const setSim = (buf: rt.Buffer, s: Sim): rt.Err => rt.setStruct(simMeta, buf, s);
export const readSim = (buf: rt.Buffer): [Sim, rt.Err] => getSim(buf);
export const eqSim = (a: Sim | null | undefined, b: Sim | null | undefined): boolean => rt.eqStruct(simMeta, a, b);

export const getSimListBody = (buf: rt.Buffer, state: number): [Sim[], rt.Err] => {
    const [list, err] = rt.getDefaultList<Sim>(
        buf,
        state,
        () => newSim(),
        (buf) => readSim(buf),
    );
    return [list, err];
};

export const setSimListBody = (buf: rt.Buffer, state: number, v: Sim[]): rt.Err => {
    return rt.setDefaultList<Sim>(
        buf,
        state,
        v,
        (item) => isZeroSim(item),
        (buf, item) => setSim(buf, item),
    );
}
