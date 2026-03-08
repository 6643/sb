import * as rt from "./type"
import * as _ from "./_"

export interface Sim extends rt.Serializable, rt.Deserializable {
    // SIM卡ID
    id: number;
    type: _.Type;
    status: _.ItemStatus;
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
    operator: _.SimOperator;
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
    pickPhone: _.SimPickPhone[];
    // 首充渠道
    firstChargeLink: string;
    // 首充金额
    firstChargeMoney: string;
    // 首充返额
    firstChargeReturn: string;
    // 禁发区域
    banCity: number[];
    info: _.SimInfo[];
    // 套餐截图
    snapshot: string[];
}

const simHeaderWidths = [1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2] as const;

export const newSim = (): Sim => {
    const s = {
        id: 0,
        type: _.DefaultType(),
        status: _.DefaultItemStatus(),
        commission: 0,
        supplier: 0,
        aff: 0,
        contractDuration: 0,
        name: "",
        operator: _.DefaultSimOperator(),
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
    } as any as Sim;
    s.set = (buf: rt.Buffer) => setSim(buf, s);
    s.get = (buf: rt.Buffer) => {
        const [next, err] = getSim(buf);
        if (err === undefined) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroSim = (s: Sim | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (!_.IsDefaultType(s.type as any)) return false;
    if (!_.IsDefaultItemStatus(s.status as any)) return false;
    if (s.commission !== 0) return false;
    if (s.supplier !== 0) return false;
    if (s.aff !== 0) return false;
    if (s.contractDuration !== 0) return false;
    if (s.name !== "") return false;
    if (!_.IsDefaultSimOperator(s.operator as any)) return false;
    if (s.monthly !== 0) return false;
    if (s.flowUniversal !== 0) return false;
    if (s.flowDirectional !== 0) return false;
    if (s.canMoveFlow) return false;
    if (s.callMonth !== 0) return false;
    if (s.callPrice !== 0) return false;
    if (s.smsMonth !== 0) return false;
    if (s.smsPrice !== 0) return false;
    if (s.minAge !== 0) return false;
    if (s.maxAge !== 0) return false;
    if (s.attribution !== 0) return false;
    if (s.pickPhone.length !== 0) return false;
    if (s.firstChargeLink !== "") return false;
    if (s.firstChargeMoney !== "") return false;
    if (s.firstChargeReturn !== "") return false;
    if (s.banCity.length !== 0) return false;
    if (s.info.length !== 0) return false;
    if (s.snapshot.length !== 0) return false;
    return true;
};

export const validateSim = (s: Sim | null | undefined): rt.Err => {
    if (s === null || s === undefined) return undefined;
    if (!_.IsAssignableType(s.type as any)) return new Error(`Type 非法枚举值: ${s.type as any}`);
    if (!_.IsAssignableItemStatus(s.status as any)) return new Error(`Status 非法枚举值: ${s.status as any}`);
    { const [, err] = rt.textState(s.name); if (err !== null) return new Error(`Name: ${err.message}`); }
    if (!_.IsAssignableSimOperator(s.operator as any)) return new Error(`Operator 非法枚举值: ${s.operator as any}`);
    for (let i = 0; i < s.pickPhone.length; i++) {
        if (!_.IsAssignableSimPickPhone(s.pickPhone[i] as any)) return new Error(`PickPhone[${i}] 非法枚举值: ${s.pickPhone[i] as any}`);
    }
    { const [, err] = rt.listCountState(s.pickPhone.length); if (err !== null) return new Error(`PickPhone: ${err.message}`); }
    { const [, err] = rt.textState(s.firstChargeLink); if (err !== null) return new Error(`FirstChargeLink: ${err.message}`); }
    { const [, err] = rt.textState(s.firstChargeMoney); if (err !== null) return new Error(`FirstChargeMoney: ${err.message}`); }
    { const [, err] = rt.textState(s.firstChargeReturn); if (err !== null) return new Error(`FirstChargeReturn: ${err.message}`); }
    { const [, err] = rt.listCountState(s.banCity.length); if (err !== null) return new Error(`BanCity: ${err.message}`); }
    { const [, err] = rt.listCountState(s.info.length); if (err !== null) return new Error(`Info: ${err.message}`); }
    for (let i = 0; i < s.info.length; i++) {
        const err = _.validateSimInfo(s.info[i]);
        if (err !== undefined) return new Error(`Info[${i}]: ${err.message}`);
    }
    { const [, err] = rt.listCountState(s.snapshot.length); if (err !== null) return new Error(`Snapshot: ${err.message}`); }
    for (let i = 0; i < s.snapshot.length; i++) {
        const [, err] = rt.textState(s.snapshot[i]);
        if (err !== null) return new Error(`Snapshot[${i}]: ${err.message}`);
    }
    return undefined;
};

export const getSim = (buf: rt.Buffer): [Sim, rt.Err] => {
    const s = newSim();
    const headerBits = 35;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const [headerStates, errHeaderState] = rt.readHeader(header, simHeaderWidths, "Sim header");
    if (errHeaderState !== undefined) return [s, errHeaderState];
    const idPresent = headerStates[0] === 1;
    const typePresent = headerStates[1] === 1;
    const statusPresent = headerStates[2] === 1;
    const commissionPresent = headerStates[3] === 1;
    const supplierPresent = headerStates[4] === 1;
    const affPresent = headerStates[5] === 1;
    const contractDurationPresent = headerStates[6] === 1;
    const nameState = headerStates[7];
    const operatorPresent = headerStates[8] === 1;
    const monthlyPresent = headerStates[9] === 1;
    const flowUniversalPresent = headerStates[10] === 1;
    const flowDirectionalPresent = headerStates[11] === 1;
    const canMoveFlowState = headerStates[12] === 1;
    const callMonthPresent = headerStates[13] === 1;
    const callPricePresent = headerStates[14] === 1;
    const smsMonthPresent = headerStates[15] === 1;
    const smsPricePresent = headerStates[16] === 1;
    const minAgePresent = headerStates[17] === 1;
    const maxAgePresent = headerStates[18] === 1;
    const attributionPresent = headerStates[19] === 1;
    const pickPhoneState = headerStates[20];
    const firstChargeLinkState = headerStates[21];
    const firstChargeMoneyState = headerStates[22];
    const firstChargeReturnState = headerStates[23];
    const banCityState = headerStates[24];
    const infoState = headerStates[25];
    const snapshotState = headerStates[26];
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get Sim Id: ${err.message}`)];
        s.id = value as any;
    }
    s.type = _.DefaultType();
    if (typePresent) {
        const [value, err] = rt.getU8(buf);
        if (err !== null) return [s, new Error(`get Sim Type: ${err.message}`)];
        const item = value as _.Type;
        if (!_.IsType(item)) return [s, new Error(`get Sim Type: 非法枚举值: ${item}`)];
        s.type = item;
    }
    s.status = _.DefaultItemStatus();
    if (statusPresent) {
        const [value, err] = rt.getU8(buf);
        if (err !== null) return [s, new Error(`get Sim Status: ${err.message}`)];
        const item = value as _.ItemStatus;
        if (!_.IsItemStatus(item)) return [s, new Error(`get Sim Status: 非法枚举值: ${item}`)];
        s.status = item;
    }
    if (commissionPresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get Sim Commission: ${err.message}`)];
        s.commission = value as any;
    }
    if (supplierPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get Sim Supplier: ${err.message}`)];
        s.supplier = value as any;
    }
    if (affPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get Sim Aff: ${err.message}`)];
        s.aff = value as any;
    }
    if (contractDurationPresent) {
        const [value, err] = rt.getU8(buf);
        if (err !== null) return [s, new Error(`get Sim ContractDuration: ${err.message}`)];
        s.contractDuration = value as any;
    }
    { const [value, err] = rt.getText(buf, nameState); if (err !== null) return [s, new Error(`get Sim Name: ${err.message}`)]; s.name = value; }
    s.operator = _.DefaultSimOperator();
    if (operatorPresent) {
        const [value, err] = rt.getU8(buf);
        if (err !== null) return [s, new Error(`get Sim Operator: ${err.message}`)];
        const item = value as _.SimOperator;
        if (!_.IsSimOperator(item)) return [s, new Error(`get Sim Operator: 非法枚举值: ${item}`)];
        s.operator = item;
    }
    if (monthlyPresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get Sim Monthly: ${err.message}`)];
        s.monthly = value as any;
    }
    if (flowUniversalPresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get Sim FlowUniversal: ${err.message}`)];
        s.flowUniversal = value as any;
    }
    if (flowDirectionalPresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get Sim FlowDirectional: ${err.message}`)];
        s.flowDirectional = value as any;
    }
    s.canMoveFlow = canMoveFlowState;
    if (callMonthPresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get Sim CallMonth: ${err.message}`)];
        s.callMonth = value as any;
    }
    if (callPricePresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get Sim CallPrice: ${err.message}`)];
        s.callPrice = value as any;
    }
    if (smsMonthPresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get Sim SmsMonth: ${err.message}`)];
        s.smsMonth = value as any;
    }
    if (smsPricePresent) {
        const [value, err] = rt.getU16(buf);
        if (err !== null) return [s, new Error(`get Sim SmsPrice: ${err.message}`)];
        s.smsPrice = value as any;
    }
    if (minAgePresent) {
        const [value, err] = rt.getU8(buf);
        if (err !== null) return [s, new Error(`get Sim MinAge: ${err.message}`)];
        s.minAge = value as any;
    }
    if (maxAgePresent) {
        const [value, err] = rt.getU8(buf);
        if (err !== null) return [s, new Error(`get Sim MaxAge: ${err.message}`)];
        s.maxAge = value as any;
    }
    if (attributionPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get Sim Attribution: ${err.message}`)];
        s.attribution = value as any;
    }
    { const [value, err] = _.getSimPickPhoneListBody(buf, pickPhoneState); if (err !== undefined) return [s, new Error(`get Sim PickPhone: ${err.message}`)]; s.pickPhone = value; }
    { const [value, err] = rt.getText(buf, firstChargeLinkState); if (err !== null) return [s, new Error(`get Sim FirstChargeLink: ${err.message}`)]; s.firstChargeLink = value; }
    { const [value, err] = rt.getText(buf, firstChargeMoneyState); if (err !== null) return [s, new Error(`get Sim FirstChargeMoney: ${err.message}`)]; s.firstChargeMoney = value; }
    { const [value, err] = rt.getText(buf, firstChargeReturnState); if (err !== null) return [s, new Error(`get Sim FirstChargeReturn: ${err.message}`)]; s.firstChargeReturn = value; }
    { const [value, err] = rt.getZeroList<number>(buf, banCityState, 0, rt.getU32); if (err !== undefined) return [s, new Error(`get Sim BanCity: ${err.message}`)]; s.banCity = value; }
    { const [value, err] = _.getSimInfoListBody(buf, infoState); if (err !== undefined) return [s, new Error(`get Sim Info: ${err.message}`)]; s.info = value; }
    { const [value, err] = rt.getTextList(buf, snapshotState); if (err !== null) return [s, new Error(`get Sim Snapshot: ${err.message}`)]; s.snapshot = value; }
    const errValidate = validateSim(s);
    if (errValidate !== undefined) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, undefined];
};

export const setSim = (buf: rt.Buffer, s: Sim): rt.Err => {
    if (s === null || s === undefined) return new Error(`set Sim: value is null or undefined`);
    const errValidate = validateSim(s);
    if (errValidate !== undefined) return new Error(`validate Sim: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [nameState, errNameState] = rt.textState(s.name);
    if (errNameState !== null) return errNameState;
    const [pickPhoneState, errPickPhoneState] = rt.listCountState(s.pickPhone.length);
    if (errPickPhoneState !== null) return errPickPhoneState;
    const [firstChargeLinkState, errFirstChargeLinkState] = rt.textState(s.firstChargeLink);
    if (errFirstChargeLinkState !== null) return errFirstChargeLinkState;
    const [firstChargeMoneyState, errFirstChargeMoneyState] = rt.textState(s.firstChargeMoney);
    if (errFirstChargeMoneyState !== null) return errFirstChargeMoneyState;
    const [firstChargeReturnState, errFirstChargeReturnState] = rt.textState(s.firstChargeReturn);
    if (errFirstChargeReturnState !== null) return errFirstChargeReturnState;
    const [banCityState, errBanCityState] = rt.listCountState(s.banCity.length);
    if (errBanCityState !== null) return errBanCityState;
    const [infoState, errInfoState] = rt.listCountState(s.info.length);
    if (errInfoState !== null) return errInfoState;
    const [snapshotState, errSnapshotState] = rt.listCountState(s.snapshot.length);
    if (errSnapshotState !== null) return errSnapshotState;
    const headerStates = [
        s.id !== 0 ? 1 : 0,
        !_.IsDefaultType(s.type as any) ? 1 : 0,
        !_.IsDefaultItemStatus(s.status as any) ? 1 : 0,
        s.commission !== 0 ? 1 : 0,
        s.supplier !== 0 ? 1 : 0,
        s.aff !== 0 ? 1 : 0,
        s.contractDuration !== 0 ? 1 : 0,
        nameState,
        !_.IsDefaultSimOperator(s.operator as any) ? 1 : 0,
        s.monthly !== 0 ? 1 : 0,
        s.flowUniversal !== 0 ? 1 : 0,
        s.flowDirectional !== 0 ? 1 : 0,
        s.canMoveFlow ? 1 : 0,
        s.callMonth !== 0 ? 1 : 0,
        s.callPrice !== 0 ? 1 : 0,
        s.smsMonth !== 0 ? 1 : 0,
        s.smsPrice !== 0 ? 1 : 0,
        s.minAge !== 0 ? 1 : 0,
        s.maxAge !== 0 ? 1 : 0,
        s.attribution !== 0 ? 1 : 0,
        pickPhoneState,
        firstChargeLinkState,
        firstChargeMoneyState,
        firstChargeReturnState,
        banCityState,
        infoState,
        snapshotState,
    ];
    const [header, errHeader] = rt.writeHeader(simHeaderWidths, headerStates);
    if (errHeader !== undefined) { buf.rewindWrite(startOffset); return new Error(`set header: ${errHeader.message}`); }
    const errHeaderWrite = buf.write(header);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Id: ${err.message}`); }
    }
    if (!_.IsDefaultType(s.type as any)) {
        const err = rt.setU8(buf, _.NormalizeType(s.type as any) as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Type: ${err.message}`); }
    }
    if (!_.IsDefaultItemStatus(s.status as any)) {
        const err = rt.setU8(buf, _.NormalizeItemStatus(s.status as any) as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Status: ${err.message}`); }
    }
    if (s.commission !== 0) {
        const err = rt.setU16(buf, s.commission as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Commission: ${err.message}`); }
    }
    if (s.supplier !== 0) {
        const err = rt.setU32(buf, s.supplier as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Supplier: ${err.message}`); }
    }
    if (s.aff !== 0) {
        const err = rt.setU32(buf, s.aff as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Aff: ${err.message}`); }
    }
    if (s.contractDuration !== 0) {
        const err = rt.setU8(buf, s.contractDuration as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim ContractDuration: ${err.message}`); }
    }
    { const err = rt.setText(buf, nameState, s.name); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Name: ${err.message}`); } }
    if (!_.IsDefaultSimOperator(s.operator as any)) {
        const err = rt.setU8(buf, _.NormalizeSimOperator(s.operator as any) as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Operator: ${err.message}`); }
    }
    if (s.monthly !== 0) {
        const err = rt.setU16(buf, s.monthly as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Monthly: ${err.message}`); }
    }
    if (s.flowUniversal !== 0) {
        const err = rt.setU16(buf, s.flowUniversal as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim FlowUniversal: ${err.message}`); }
    }
    if (s.flowDirectional !== 0) {
        const err = rt.setU16(buf, s.flowDirectional as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim FlowDirectional: ${err.message}`); }
    }
    if (s.callMonth !== 0) {
        const err = rt.setU16(buf, s.callMonth as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim CallMonth: ${err.message}`); }
    }
    if (s.callPrice !== 0) {
        const err = rt.setU16(buf, s.callPrice as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim CallPrice: ${err.message}`); }
    }
    if (s.smsMonth !== 0) {
        const err = rt.setU16(buf, s.smsMonth as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim SmsMonth: ${err.message}`); }
    }
    if (s.smsPrice !== 0) {
        const err = rt.setU16(buf, s.smsPrice as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim SmsPrice: ${err.message}`); }
    }
    if (s.minAge !== 0) {
        const err = rt.setU8(buf, s.minAge as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim MinAge: ${err.message}`); }
    }
    if (s.maxAge !== 0) {
        const err = rt.setU8(buf, s.maxAge as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim MaxAge: ${err.message}`); }
    }
    if (s.attribution !== 0) {
        const err = rt.setU32(buf, s.attribution as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Attribution: ${err.message}`); }
    }
    { const err = _.setSimPickPhoneListBody(buf, pickPhoneState, s.pickPhone); if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set Sim PickPhone: ${err.message}`); } }
    { const err = rt.setText(buf, firstChargeLinkState, s.firstChargeLink); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim FirstChargeLink: ${err.message}`); } }
    { const err = rt.setText(buf, firstChargeMoneyState, s.firstChargeMoney); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim FirstChargeMoney: ${err.message}`); } }
    { const err = rt.setText(buf, firstChargeReturnState, s.firstChargeReturn); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim FirstChargeReturn: ${err.message}`); } }
    { const err = rt.setZeroList<number>(buf, banCityState, s.banCity, 0, rt.setU32); if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set Sim BanCity: ${err.message}`); } }
    { const err = _.setSimInfoListBody(buf, infoState, s.info); if (err !== undefined) { buf.rewindWrite(startOffset); return new Error(`set Sim Info: ${err.message}`); } }
    { const err = rt.setTextList(buf, snapshotState, s.snapshot); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set Sim Snapshot: ${err.message}`); } }
    return undefined;
};

export const readSim = (buf: rt.Buffer): [Sim, rt.Err] => getSim(buf);

export const eqSim = (a: Sim | null | undefined, b: Sim | null | undefined): boolean => {
    if (isZeroSim(a as any) && isZeroSim(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!_.eqTypeValue(a.type as any, b.type as any)) return false;
    if (!_.eqItemStatusValue(a.status as any, b.status as any)) return false;
    if (!rt.eqU16(a.commission, b.commission)) return false;
    if (!rt.eqU32(a.supplier, b.supplier)) return false;
    if (!rt.eqU32(a.aff, b.aff)) return false;
    if (!rt.eqU8(a.contractDuration, b.contractDuration)) return false;
    if (!rt.eqText(a.name, b.name)) return false;
    if (!_.eqSimOperatorValue(a.operator as any, b.operator as any)) return false;
    if (!rt.eqU16(a.monthly, b.monthly)) return false;
    if (!rt.eqU16(a.flowUniversal, b.flowUniversal)) return false;
    if (!rt.eqU16(a.flowDirectional, b.flowDirectional)) return false;
    if (!rt.eqBool(a.canMoveFlow, b.canMoveFlow)) return false;
    if (!rt.eqU16(a.callMonth, b.callMonth)) return false;
    if (!rt.eqU16(a.callPrice, b.callPrice)) return false;
    if (!rt.eqU16(a.smsMonth, b.smsMonth)) return false;
    if (!rt.eqU16(a.smsPrice, b.smsPrice)) return false;
    if (!rt.eqU8(a.minAge, b.minAge)) return false;
    if (!rt.eqU8(a.maxAge, b.maxAge)) return false;
    if (!rt.eqU32(a.attribution, b.attribution)) return false;
    if (!_.eqSimPickPhoneList(a.pickPhone as any, b.pickPhone as any)) return false;
    if (!rt.eqText(a.firstChargeLink, b.firstChargeLink)) return false;
    if (!rt.eqText(a.firstChargeMoney, b.firstChargeMoney)) return false;
    if (!rt.eqText(a.firstChargeReturn, b.firstChargeReturn)) return false;
    if (!rt.eqList(a.banCity, b.banCity, rt.eqU32)) return false;
    if (!rt.eqList(a.info, b.info, _.eqSimInfo)) return false;
    if (!rt.eqList(a.snapshot, b.snapshot, rt.eqText)) return false;
    return true;
};

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
