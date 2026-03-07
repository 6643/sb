import * as _ from "./_"
import * as TplEnum from "./enum"

export interface RechargeA extends _.Serializable, _.Deserializable {
    id: number;
    type: TplEnum.OrderStatus[];
    phone: string[];
    si: _.SimInfo | null;
    aid: number;
}

export const newRechargeA = (): RechargeA => {
    const s = {
        id: 0,
        type: [],
        phone: [],
        si: null,
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
    if ((a.si === null) !== (b.si === null)) return false;
    if (a.si !== null && b.si !== null && !_.eqSimInfo(a.si, b.si)) return false;
    if (!_.eqU32(a.aid, b.aid)) return false;
    return true;
}

export const getRechargeA = (buf: _.Buffer): [RechargeA, Error | null] => {
    const s = newRechargeA();
    const bitmaskSize = 1;
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
        if (!_.IsOrderStatusList(s.type as any)) return [s, new Error("get RechargeA Type: invalid enum value")];
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
    const startOffset = buf.write_offset;
    const bits = new Uint8Array(1);
    if (!_.eqU32(s.id, 0)) {
        _.SetBit(bits, 0, true);
    }
    if (s.type && s.type.length > 0) {
        if (!_.IsOrderStatusList(s.type as any)) return new Error("set RechargeA Type: invalid enum value");
        _.SetBit(bits, 1, true);
    }
    if (s.phone && s.phone.length > 0) {
        _.SetBit(bits, 2, true);
    }
    if (s.si !== null && s.si !== undefined) {
        _.SetBit(bits, 3, true);
    }
    if (!_.eqU32(s.aid, 0)) {
        _.SetBit(bits, 4, true);
    }

    const errBits = buf.write(bits);
    if (errBits !== null) return errBits;
    if (!_.eqU32(s.id, 0)) {
        const err = _.setU32(buf, s.id);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    if (s.type && s.type.length > 0) {
        const err = _.setU8List(buf, s.type as any);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    if (s.phone && s.phone.length > 0) {
        const err = _.setTextList(buf, s.phone);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    if (s.si !== null && s.si !== undefined) {
        const err = _.setSimInfo(buf, s.si);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    if (!_.eqU32(s.aid, 0)) {
        const err = _.setU32(buf, s.aid);
        if (err !== null) {
            buf.rewindWrite(startOffset);
            return err;
        }
    }
    return null;
}

export const getRechargeAList = (buf: _.Buffer): [RechargeA[], Error | null] => {
    const [count, err] = _.getU16(buf);
    if (err !== null) return [[], err];
    const list: RechargeA[] = new Array(count);
    for (let i = 0; i < count; i++) {
        const [item, err2] = getRechargeA(buf);
        if (err2 !== null) return [[], err2];
        list[i] = item;
    }
    return [list, null];
}
export const setRechargeAList = (buf: _.Buffer, v: RechargeA[]): Error | null => {
    if (v.length > 65535) return new Error(`list length ${v.length} exceeds u16 max`);
    const err = _.setU16(buf, v.length);
    if (err !== null) return err;
    for (const item of v) {
        const err2 = setRechargeA(buf, item);
        if (err2 !== null) return err2;
    }
    return null;
}
export const eqRechargeAList = (a: RechargeA[], b: RechargeA[]): boolean => _.eqList(a, b, eqRechargeA);
