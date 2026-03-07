import * as rt from "./type"
import * as _ from "./_"

export interface RechargeB extends rt.Serializable, rt.Deserializable {
    // abcd
    id: number;
    type: _.OrderStatus[];
    phone: string[];
    si: _.SimInfo | null;
    bid: number;
}

export const newRechargeB = (): RechargeB => {
    const s = {
        id: 0,
        type: [],
        phone: [],
        si: null,
        bid: 0,
    } as any as RechargeB;
    s.set = (buf: rt.Buffer) => setRechargeB(buf, s);
    s.get = (buf: rt.Buffer) => {
        const [next, err] = getRechargeB(buf);
        if (err === null) Object.assign(s, next);
        return err;
    };
    return s;
};

export const isZeroRechargeB = (s: RechargeB | null | undefined): boolean => {
    if (s === null || s === undefined) return true;
    if (s.id !== 0) return false;
    if (s.type.length !== 0) return false;
    if (s.phone.length !== 0) return false;
    if (!_.isZeroSimInfo(s.si)) return false;
    if (s.bid !== 0) return false;
    return true;
};

export const validateRechargeB = (s: RechargeB | null | undefined): Error | null => {
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

export const getRechargeB = (buf: rt.Buffer): [RechargeB, Error | null] => {
    const s = newRechargeB();
    const headerBits = 7;
    const [header, errHeader] = buf.read(rt.headerSize(headerBits));
    if (errHeader !== null) return [s, new Error(`not enough data`)];
    const errPadding = rt.validatePaddingZero(header, headerBits, "RechargeB header");
    if (errPadding !== null) return [s, errPadding];
    const reader = new rt.BitReader(header, headerBits);
    const [idPresent, errIdPresent] = reader.readBit();
    if (errIdPresent !== null) return [s, new Error(`get RechargeB Id header: ${errIdPresent.message}`)];
    const [typeState, errTypeState] = reader.readBits(2);
    if (errTypeState !== null) return [s, new Error(`get RechargeB Type header: ${errTypeState.message}`)];
    const [phoneState, errPhoneState] = reader.readBits(2);
    if (errPhoneState !== null) return [s, new Error(`get RechargeB Phone header: ${errPhoneState.message}`)];
    const [siPresent, errSiPresent] = reader.readBit();
    if (errSiPresent !== null) return [s, new Error(`get RechargeB Si header: ${errSiPresent.message}`)];
    const [bidPresent, errBidPresent] = reader.readBit();
    if (errBidPresent !== null) return [s, new Error(`get RechargeB Bid header: ${errBidPresent.message}`)];
    if (idPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get RechargeB Id: ${err.message}`)];
        s.id = value as any;
    }
    { const [value, err] = _.getOrderStatusListBody(buf, typeState); if (err !== null) return [s, new Error(`get RechargeB Type: ${err.message}`)]; s.type = value; }
    { const [value, err] = rt.getTextListCompact(buf, phoneState); if (err !== null) return [s, new Error(`get RechargeB Phone: ${err.message}`)]; s.phone = value; }
    if (siPresent) {
        const [value, err] = _.readSimInfo(buf);
        if (err !== null) return [s, new Error(`get RechargeB Si: ${err.message}`)];
        s.si = value;
    }
    if (bidPresent) {
        const [value, err] = rt.getU32(buf);
        if (err !== null) return [s, new Error(`get RechargeB Bid: ${err.message}`)];
        s.bid = value as any;
    }
    const errValidate = validateRechargeB(s);
    if (errValidate !== null) return [s, new Error(`validate failed: ${errValidate.message}`)];
    return [s, null];
};

export const setRechargeB = (buf: rt.Buffer, s: RechargeB): Error | null => {
    if (s === null || s === undefined) return new Error(`set RechargeB: value is null or undefined`);
    const errValidate = validateRechargeB(s);
    if (errValidate !== null) return new Error(`validate RechargeB: ${errValidate.message}`);
    const startOffset = buf.writeOffset;
    const [typeState, errTypeState] = rt.listCountState(s.type.length);
    if (errTypeState !== null) return errTypeState;
    const [phoneState, errPhoneState] = rt.listCountState(s.phone.length);
    if (errPhoneState !== null) return errPhoneState;
    const header = new rt.BitWriter(7);
    header.writeBit(s.id !== 0);
    { const err = header.writeBits(typeState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeB Type header: ${err.message}`); } }
    { const err = header.writeBits(phoneState, 2); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeB Phone header: ${err.message}`); } }
    header.writeBit(!_.isZeroSimInfo(s.si));
    header.writeBit(s.bid !== 0);
    const errHeaderWrite = buf.write(header.bytes);
    if (errHeaderWrite !== null) { buf.rewindWrite(startOffset); return errHeaderWrite; }
    if (s.id !== 0) {
        const err = rt.setU32(buf, s.id as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeB Id: ${err.message}`); }
    }
    { const err = _.setOrderStatusListBody(buf, typeState, s.type); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeB Type: ${err.message}`); } }
    { const err = rt.setTextListCompact(buf, phoneState, s.phone); if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeB Phone: ${err.message}`); } }
    if (!_.isZeroSimInfo(s.si)) {
        const err = _.setSimInfo(buf, s.si!);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeB Si: ${err.message}`); }
    }
    if (s.bid !== 0) {
        const err = rt.setU32(buf, s.bid as any);
        if (err !== null) { buf.rewindWrite(startOffset); return new Error(`set RechargeB Bid: ${err.message}`); }
    }
    return null;
};

export const readRechargeB = (buf: rt.Buffer): [RechargeB, Error | null] => getRechargeB(buf);

export const eqRechargeB = (a: RechargeB | null | undefined, b: RechargeB | null | undefined): boolean => {
    if (isZeroRechargeB(a as any) && isZeroRechargeB(b as any)) return true;
    if (a === null || a === undefined || b === null || b === undefined) return false;
    if (!rt.eqU32(a.id, b.id)) return false;
    if (!_.eqOrderStatusList(a.type as any, b.type as any)) return false;
    if (!rt.eqList(a.phone, b.phone, rt.eqText)) return false;
    if (!_.eqSimInfo(a.si, b.si)) return false;
    if (!rt.eqU32(a.bid, b.bid)) return false;
    return true;
};

export const getRechargeBListBody = (buf: rt.Buffer, state: number): [RechargeB[], Error | null] =>
    rt.getBitmapListCompact<RechargeB>(
        buf,
        state,
        () => newRechargeB(),
        (buf) => readRechargeB(buf),
    );

export const setRechargeBListBody = (buf: rt.Buffer, state: number, v: RechargeB[]): Error | null =>
    rt.setBitmapListCompact<RechargeB>(
        buf,
        state,
        v,
        (item) => isZeroRechargeB(item),
        (buf, item) => setRechargeB(buf, item),
    );
