import * as _ from "./_.ts"
import * as Enum from "./enum"

export interface Sim extends _.Serializable, _.Deserializable {
    id: number;
    type: Enum.Type;
    status: Enum.ItemStatus;
    commission: number;
    supplier: number;
    aff: number;
    contractDuration: number;
    name: string;
    operator: Enum.SimOperator;
    monthly: number;
    flowUniversal: number;
    flowDirectional: number;
    canMoveFlow: boolean;
    callMonth: number;
    callPrice: number;
    smsMonth: number;
    smsPrice: number;
    minAge: number;
    maxAge: number;
    attribution: number;
    pickPhone: Enum.SimPickPhone[];
    firstChargeLink: string;
    firstChargeMoney: string;
    firstChargeReturn: string;
    banCity: number[];
    info: _.SimInfo[];
    snapshot: string[];
}

export const newSim = (): Sim => {
    const s = {
        id: 0,
        type: 0,
        status: 0,
        commission: 0,
        supplier: 0,
        aff: 0,
        contractDuration: 0,
        name: "",
        operator: 0,
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
    s.set = (buf: _.Buffer) => setSim(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getSim(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqSim = (a: Sim, b: Sim): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqU32(a.id, b.id)) return false;
    if (a.type !== b.type) return false;
    if (a.status !== b.status) return false;
    if (!_.eqU16(a.commission, b.commission)) return false;
    if (!_.eqU32(a.supplier, b.supplier)) return false;
    if (!_.eqU32(a.aff, b.aff)) return false;
    if (!_.eqU8(a.contractDuration, b.contractDuration)) return false;
    if (!_.eqText(a.name, b.name)) return false;
    if (a.operator !== b.operator) return false;
    if (!_.eqU16(a.monthly, b.monthly)) return false;
    if (!_.eqU16(a.flowUniversal, b.flowUniversal)) return false;
    if (!_.eqU16(a.flowDirectional, b.flowDirectional)) return false;
    if (!_.eqBool(a.canMoveFlow, b.canMoveFlow)) return false;
    if (!_.eqU16(a.callMonth, b.callMonth)) return false;
    if (!_.eqU16(a.callPrice, b.callPrice)) return false;
    if (!_.eqU16(a.smsMonth, b.smsMonth)) return false;
    if (!_.eqU16(a.smsPrice, b.smsPrice)) return false;
    if (!_.eqU8(a.minAge, b.minAge)) return false;
    if (!_.eqU8(a.maxAge, b.maxAge)) return false;
    if (!_.eqU32(a.attribution, b.attribution)) return false;
    if (!_.eqU8List(a.pickPhone as any, b.pickPhone as any)) return false;
    if (!_.eqText(a.firstChargeLink, b.firstChargeLink)) return false;
    if (!_.eqText(a.firstChargeMoney, b.firstChargeMoney)) return false;
    if (!_.eqText(a.firstChargeReturn, b.firstChargeReturn)) return false;
    if (!_.eqU32List(a.banCity, b.banCity)) return false;
    if (!_.eqSimInfoList(a.info, b.info)) return false;
    if (!_.eqTextList(a.snapshot, b.snapshot)) return false;
    return true;
}

export const getSim = (buf: _.Buffer): [Sim, Error | null] => {
    const s = newSim();
    const bitmaskSize = Math.ceil(27 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.id = v;
    }
    if (_.GetBit(bits, 1)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.type = v as any;
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.status = v as any;
    }
    if (_.GetBit(bits, 3)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.commission = v;
    }
    if (_.GetBit(bits, 4)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.supplier = v;
    }
    if (_.GetBit(bits, 5)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.aff = v;
    }
    if (_.GetBit(bits, 6)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.contractDuration = v;
    }
    if (_.GetBit(bits, 7)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.name = v;
    }
    if (_.GetBit(bits, 8)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.operator = v as any;
    }
    if (_.GetBit(bits, 9)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.monthly = v;
    }
    if (_.GetBit(bits, 10)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.flowUniversal = v;
    }
    if (_.GetBit(bits, 11)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.flowDirectional = v;
    }
    s.canMoveFlow = _.GetBit(bits, 12);
    if (_.GetBit(bits, 13)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.callMonth = v;
    }
    if (_.GetBit(bits, 14)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.callPrice = v;
    }
    if (_.GetBit(bits, 15)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.smsMonth = v;
    }
    if (_.GetBit(bits, 16)) {
        const [v, err] = _.getU16(buf);
        if (err !== null) return [s, err];
        s.smsPrice = v;
    }
    if (_.GetBit(bits, 17)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.minAge = v;
    }
    if (_.GetBit(bits, 18)) {
        const [v, err] = _.getU8(buf);
        if (err !== null) return [s, err];
        s.maxAge = v;
    }
    if (_.GetBit(bits, 19)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.attribution = v;
    }
    if (_.GetBit(bits, 20)) {
        const [v, err] = _.getU8List(buf);
        if (err !== null) return [s, err];
        s.pickPhone = v as any;
    }
    if (_.GetBit(bits, 21)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.firstChargeLink = v;
    }
    if (_.GetBit(bits, 22)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.firstChargeMoney = v;
    }
    if (_.GetBit(bits, 23)) {
        const [v, err] = _.getText(buf);
        if (err !== null) return [s, err];
        s.firstChargeReturn = v;
    }
    if (_.GetBit(bits, 24)) {
        const [v, err] = _.getU32List(buf);
        if (err !== null) return [s, err];
        s.banCity = v;
    }
    if (_.GetBit(bits, 25)) {
        const [v, err] = _.getSimInfoList(buf);
        if (err !== null) return [s, err];
        s.info = v;
    }
    if (_.GetBit(bits, 26)) {
        const [v, err] = _.getTextList(buf);
        if (err !== null) return [s, err];
        s.snapshot = v;
    }
    return [s, null];
}

export const setSim = (buf: _.Buffer, s: Sim): Error | null => {
    if (s === null || s === undefined) return new Error(`set Sim: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(27 / 8));
    const body = new _.Buffer();
    if (!_.eqU32(s.id, 0)) {
        const err = _.setU32(body, s.id);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }
    if ((s.type as any) !== 0) {
        const err = _.setU8(body, s.type as any);
        if (err !== null) return err;
        _.SetBit(bits, 1, true);
    }
    if ((s.status as any) !== 0) {
        const err = _.setU8(body, s.status as any);
        if (err !== null) return err;
        _.SetBit(bits, 2, true);
    }
    if (!_.eqU16(s.commission, 0)) {
        const err = _.setU16(body, s.commission);
        if (err !== null) return err;
        _.SetBit(bits, 3, true);
    }
    if (!_.eqU32(s.supplier, 0)) {
        const err = _.setU32(body, s.supplier);
        if (err !== null) return err;
        _.SetBit(bits, 4, true);
    }
    if (!_.eqU32(s.aff, 0)) {
        const err = _.setU32(body, s.aff);
        if (err !== null) return err;
        _.SetBit(bits, 5, true);
    }
    if (!_.eqU8(s.contractDuration, 0)) {
        const err = _.setU8(body, s.contractDuration);
        if (err !== null) return err;
        _.SetBit(bits, 6, true);
    }
    if (!_.eqText(s.name, "")) {
        const err = _.setText(body, s.name);
        if (err !== null) return err;
        _.SetBit(bits, 7, true);
    }
    if ((s.operator as any) !== 0) {
        const err = _.setU8(body, s.operator as any);
        if (err !== null) return err;
        _.SetBit(bits, 8, true);
    }
    if (!_.eqU16(s.monthly, 0)) {
        const err = _.setU16(body, s.monthly);
        if (err !== null) return err;
        _.SetBit(bits, 9, true);
    }
    if (!_.eqU16(s.flowUniversal, 0)) {
        const err = _.setU16(body, s.flowUniversal);
        if (err !== null) return err;
        _.SetBit(bits, 10, true);
    }
    if (!_.eqU16(s.flowDirectional, 0)) {
        const err = _.setU16(body, s.flowDirectional);
        if (err !== null) return err;
        _.SetBit(bits, 11, true);
    }
    _.SetBit(bits, 12, s.canMoveFlow as boolean);
    if (!_.eqU16(s.callMonth, 0)) {
        const err = _.setU16(body, s.callMonth);
        if (err !== null) return err;
        _.SetBit(bits, 13, true);
    }
    if (!_.eqU16(s.callPrice, 0)) {
        const err = _.setU16(body, s.callPrice);
        if (err !== null) return err;
        _.SetBit(bits, 14, true);
    }
    if (!_.eqU16(s.smsMonth, 0)) {
        const err = _.setU16(body, s.smsMonth);
        if (err !== null) return err;
        _.SetBit(bits, 15, true);
    }
    if (!_.eqU16(s.smsPrice, 0)) {
        const err = _.setU16(body, s.smsPrice);
        if (err !== null) return err;
        _.SetBit(bits, 16, true);
    }
    if (!_.eqU8(s.minAge, 0)) {
        const err = _.setU8(body, s.minAge);
        if (err !== null) return err;
        _.SetBit(bits, 17, true);
    }
    if (!_.eqU8(s.maxAge, 0)) {
        const err = _.setU8(body, s.maxAge);
        if (err !== null) return err;
        _.SetBit(bits, 18, true);
    }
    if (!_.eqU32(s.attribution, 0)) {
        const err = _.setU32(body, s.attribution);
        if (err !== null) return err;
        _.SetBit(bits, 19, true);
    }
    if (s.pickPhone && s.pickPhone.length > 0) {
        const err = _.setU8List(body, s.pickPhone as any);
        if (err !== null) return err;
        _.SetBit(bits, 20, true);
    }
    if (!_.eqText(s.firstChargeLink, "")) {
        const err = _.setText(body, s.firstChargeLink);
        if (err !== null) return err;
        _.SetBit(bits, 21, true);
    }
    if (!_.eqText(s.firstChargeMoney, "")) {
        const err = _.setText(body, s.firstChargeMoney);
        if (err !== null) return err;
        _.SetBit(bits, 22, true);
    }
    if (!_.eqText(s.firstChargeReturn, "")) {
        const err = _.setText(body, s.firstChargeReturn);
        if (err !== null) return err;
        _.SetBit(bits, 23, true);
    }
    if (s.banCity && s.banCity.length > 0) {
        const err = _.setU32List(body, s.banCity);
        if (err !== null) return err;
        _.SetBit(bits, 24, true);
    }
    if (s.info && s.info.length > 0) {
        const err = _.setSimInfoList(body, s.info);
        if (err !== null) return err;
        _.SetBit(bits, 25, true);
    }
    if (s.snapshot && s.snapshot.length > 0) {
        const err = _.setTextList(body, s.snapshot);
        if (err !== null) return err;
        _.SetBit(bits, 26, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getSimList = (buf: _.Buffer): [Sim[], Error | null] => _.getList(buf, getSim);
export const setSimList = (buf: _.Buffer, v: Sim[]): Error | null => _.setList(buf, v, setSim);
export const eqSimList = (a: Sim[], b: Sim[]): boolean => _.eqList(a, b, eqSim);