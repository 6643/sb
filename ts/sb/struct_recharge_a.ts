import * as _ from "./_.ts"
import * as Enum from "./enum"

export interface RechargeA extends _.Serializable, _.Deserializable {
    id: number;
    type: Enum.OrderStatus[];
    phone: string[];
    si: _.SimInfo;
    aid: number;
}

export const newRechargeA = (): RechargeA => {
    const s = {
        id: 0,
        type: [],
        phone: [],
        si: _.newSimInfo(),
        aid: 0,
    } as any as RechargeA;
    s.set = (buf: _.Buffer) => setRechargeA(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getRechargeA(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqRechargeA = (a: RechargeA, b: RechargeA): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqU32(a.id, b.id)) return false;
    if (!_.eqU8List(a.type as any, b.type as any)) return false;
    if (!_.eqTextList(a.phone, b.phone)) return false;
    if (!_.eqSimInfo(a.si, b.si)) return false;
    if (!_.eqU32(a.aid, b.aid)) return false;
    return true;
}

export const getRechargeA = (buf: _.Buffer): [RechargeA, Error | null] => {
    const s = newRechargeA();
    const bitmaskSize = Math.ceil(5 / 8);
    const [bits, err] = buf.read(bitmaskSize);
    if (err !== null) return [s, err];
    if (_.GetBit(bits, 0)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.id = v;
    }
    if (_.GetBit(bits, 1)) {
        const [v, err] = _.getU8List(buf);
        if (err !== null) return [s, err];
        s.type = v as any;
    }
    if (_.GetBit(bits, 2)) {
        const [v, err] = _.getTextList(buf);
        if (err !== null) return [s, err];
        s.phone = v;
    }
    if (_.GetBit(bits, 3)) {
        const [v, err] = _.getSimInfo(buf);
        if (err !== null) return [s, err];
        s.si = v;
    }
    if (_.GetBit(bits, 4)) {
        const [v, err] = _.getU32(buf);
        if (err !== null) return [s, err];
        s.aid = v;
    }
    return [s, null];
}

export const setRechargeA = (buf: _.Buffer, s: RechargeA): Error | null => {
    if (s === null || s === undefined) return new Error(`set RechargeA: value is null or undefined`);
    const bits = new Uint8Array(Math.ceil(5 / 8));
    const body = new _.Buffer();
    if (!_.eqU32(s.id, 0)) {
        const err = _.setU32(body, s.id);
        if (err !== null) return err;
        _.SetBit(bits, 0, true);
    }
    if (s.type && s.type.length > 0) {
        const err = _.setU8List(body, s.type as any);
        if (err !== null) return err;
        _.SetBit(bits, 1, true);
    }
    if (s.phone && s.phone.length > 0) {
        const err = _.setTextList(body, s.phone);
        if (err !== null) return err;
        _.SetBit(bits, 2, true);
    }
    if (s.si !== null) {
        const err = _.setSimInfo(body, s.si);
        if (err !== null) return err;
        _.SetBit(bits, 3, true);
    }
    if (!_.eqU32(s.aid, 0)) {
        const err = _.setU32(body, s.aid);
        if (err !== null) return err;
        _.SetBit(bits, 4, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getRechargeAList = (buf: _.Buffer): [RechargeA[], Error | null] => _.getList(buf, getRechargeA);
export const setRechargeAList = (buf: _.Buffer, v: RechargeA[]): Error | null => _.setList(buf, v, setRechargeA);
export const eqRechargeAList = (a: RechargeA[], b: RechargeA[]): boolean => _.eqList(a, b, eqRechargeA);