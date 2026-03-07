import * as rt from "./type"
import * as _ from "./_"

export interface RechargeA extends rt.Serializable, rt.Deserializable {
    // abcd
    id: number;
    type: _.OrderStatus[];
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
    s.set = (buf: rt.Buffer) => setRechargeA(buf, s);
    s.get = (buf: rt.Buffer) => {
        const [next, err] = getRechargeA(buf);
        if (err === null) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroRechargeA = (s: RechargeA | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (s.type.length !== 0) return false;
    if (s.phone.length !== 0) return false;
    if (!_.isZeroSimInfo(s.si)) return false;
    if (s.aid !== 0) return false;
    return true;
};

export const validateRechargeA = (s: RechargeA | null | undefined): Error | null => {
    if (s === null || s === undefined) return null;
    for (let i = 0; i < s.type.length; i++) {
        if (!_.IsAssignableOrderStatus(s.type[i] as any)) return new Error(`Type[${i}] 非法枚举值: ${s.type[i] as any}`);
    }
    { const [, err] = rt.listCountState(s.type.length); if (err !== null) return new Error(`Type: ${err.message}`); }
    { const [, err] = rt.listCountState(s.phone.length); if (err !== null) return new Error(`Phone: ${err.message}`); }
    for (let i = 0; i < s.phone.length; i++) {
        const [, err] = rt.textState(s.phone[i]);
        if (err !== null) return new Error(`Phone[${i}]: ${err.message}`);
    }
    { const err = _.validateSimInfo(s.si); if (err !== null) return new Error(`Si: ${err.message}`); }
    return null;
};

export const getRechargeA = (buf: rt.Buffer): [RechargeA, Error | null] => {
    const s = newRechargeA();
    const headerBits = 7;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const errPadding = rt.validatePaddingZero(header, headerBits, "RechargeA header");
    if (errPadding !== null) return [s, errPadding];
    const reader = new rt.BitReader(header, headerBits);
    const [idPresent, errIdPresent] = reader.readBit();
    if (errIdPresent !== null) return [s, new Error(`get RechargeA Id header: ${errIdPresent.message}`)];
    const [typeState, errTypeState] = reader.readBits(2);
    if (errTypeState !== null) return [s, new Error(`get RechargeA Type header: ${errTypeState.message}`)];
    const [phoneState, errPhoneState] = reader.readBits(2);
    if (errPhoneState !== null) return [s, new Error(`get RechargeA Phone header: ${errPhoneState.message}`)];
    const [siPresent, errSiPresent] = reader.readBit();
    if (errSiPresent !== null) return [s, new Error(`get RechargeA Si header: ${errSiPresent.message}`)];
    const [aidPresent, errAidPresent] = reader.readBit();
    if (errAidPresent !== null) return [s, new Error(`get RechargeA Aid header: ${errAidPresent.message}`)];
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get RechargeA Id: ${err.message}`)];
        s.id = value as any;
    }
    { const [value, err] = _.getOrderStatusListBody(buf, typeState); if (err !== null) return [s, new Error(`get RechargeA Type: ${err.message}`)]; s.type = value; }
    { const [value, err] = rt.getTextListCompact(buf, phoneState); if (err !== null) return [s, new Error(`get RechargeA Phone: ${err.message}`)]; s.phone = value; }
    if (siPresent) {
        const [value, err] = _.readSimInfo(buf);
        if (err !== null) return [s, new Error(`get RechargeA Si: ${err.message}`)];
        s.si = value;
    }
    if (aidPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get RechargeA Aid: ${err.message}`)];
        s.aid = value as any;
    }
    const errValidate = validateRechargeA(s);
    if (errValidate !== null) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, null];
};

export const setRechargeA = (buf: rt.Buffer, s: RechargeA): Error | null => {
    if (s === null || s === undefined) return new Error(`set RechargeA: value is null or undefined`);
    const errValidate = validateRechargeA(s);
    if (errValidate !== null) return new Error(`validate RechargeA: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [typeState, errTypeState] = rt.listCountState(s.type.length);
    if (errTypeState !== null) return errTypeState;
    const [phoneState, errPhoneState] = rt.listCountState(s.phone.length);
    if (errPhoneState !== null) return errPhoneState;
    const header = new rt.BitWriter(7);
    header.writeBit(s.id !== 0);
    { const err = header.writeBits(typeState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Type header: ${err.message}`); } }
    { const err = header.writeBits(phoneState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Phone header: ${err.message}`); } }
    header.writeBit(!_.isZeroSimInfo(s.si));
    header.writeBit(s.aid !== 0);
    const errHeaderWrite = buf.write(header.bytes);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Id: ${err.message}`); }
    }
    { const err = _.setOrderStatusListBody(buf, typeState, s.type); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Type: ${err.message}`); } }
    { const err = rt.setTextListCompact(buf, phoneState, s.phone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Phone: ${err.message}`); } }
    if (!_.isZeroSimInfo(s.si)) {
        const err = _.setSimInfo(buf, s.si!);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Si: ${err.message}`); }
    }
    if (s.aid !== 0) {
        const err = rt.setU32(buf, s.aid as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeA Aid: ${err.message}`); }
    }
    return null;
};

export const readRechargeA = (buf: rt.Buffer): [RechargeA, Error | null] => getRechargeA(buf);

export const eqRechargeA = (a: RechargeA | null | undefined, b: RechargeA | null | undefined): boolean => {
    if (isZeroRechargeA(a as any) && isZeroRechargeA(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!_.eqOrderStatusList(a.type as any, b.type as any)) return false;
    if (!rt.eqList(a.phone, b.phone, rt.eqText)) return false;
    if (!_.eqSimInfo(a.si, b.si)) return false;
    if (!rt.eqU32(a.aid, b.aid)) return false;
    return true;
};

export const getRechargeAListBody = (buf: rt.Buffer, state: number): [RechargeA[], Error | null] =>
    rt.getBitmapListCompact<RechargeA>(
        buf,
        state,
        () => newRechargeA(),
        (buf) => readRechargeA(buf),
    );

export const setRechargeAListBody = (buf: rt.Buffer, state: number, v: RechargeA[]): Error | null =>
    rt.setBitmapListCompact<RechargeA>(
        buf,
        state,
        v,
        (item) => isZeroRechargeA(item),
        (buf, item) => setRechargeA(buf, item),
    );
