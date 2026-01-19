import * as _ from "./_.ts"
import * as Enum from "./enum"

export interface RechargeB extends _.Serializable, _.Deserializable {
    id: number;
    type: Enum.OrderStatus[];
    phone: string[];
    si: _.SimInfo;
    bid: number;
}

export const newRechargeB = (): RechargeB => {
    const s = {
        id: 0,
        type: [],
        phone: [],
        si: _.newSimInfo(),
        bid: 0,
    } as any as RechargeB;
    s.set = (buf: _.Buffer) => setRechargeB(buf, s);
    s.get = (buf: _.Buffer) => {
        const [res, err] = getRechargeB(buf);
        if (err === null) Object.assign(s, res);
        return err;
    };
    return s;
}

export const eqRechargeB = (a: RechargeB, b: RechargeB): boolean => {
    if (a === b) return true;
    if (a === null || b === null) return false;
    if (!_.eqU32(a.id, b.id)) return false;
    if (!_.eqU8List(a.type as any, b.type as any)) return false;
    if (!_.eqTextList(a.phone, b.phone)) return false;
    if (!_.eqSimInfo(a.si, b.si)) return false;
    if (!_.eqU32(a.bid, b.bid)) return false;
    return true;
}

export const getRechargeB = (buf: _.Buffer): [RechargeB, Error | null] => {
    const s = newRechargeB();
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
        s.bid = v;
    }
    return [s, null];
}

export const setRechargeB = (buf: _.Buffer, s: RechargeB): Error | null => {
    if (s === null || s === undefined) return new Error(`set RechargeB: value is null or undefined`);
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
    if (!_.eqU32(s.bid, 0)) {
        const err = _.setU32(body, s.bid);
        if (err !== null) return err;
        _.SetBit(bits, 4, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    return buf.write(body.bytes);
}

export const getRechargeBList = (buf: _.Buffer): [RechargeB[], Error | null] => _.getList(buf, getRechargeB);
export const setRechargeBList = (buf: _.Buffer, v: RechargeB[]): Error | null => _.setList(buf, v, setRechargeB);
export const eqRechargeBList = (a: RechargeB[], b: RechargeB[]): boolean => _.eqList(a, b, eqRechargeB);